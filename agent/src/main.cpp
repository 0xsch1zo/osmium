#include "agent.hpp"

int main(int argc, char *argv[]) {
  std::unique_ptr<Agent> agent;
  if (argc == 2) {
    agent = std::make_unique<Agent>(std::stoi(argv[1]), "");
  } else {
    agent = std::make_unique<Agent>();
  }

  agent->mainLoop();
}
