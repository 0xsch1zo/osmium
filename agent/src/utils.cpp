#include "utils.hpp"
#include <array>
#include <cctype>
#include <stdexcept>
#include <windows.h>

std::string exec(const std::string &cmd) {
  auto pipe = popen((cmd + " 2>&1").c_str(), "r"); // get rid of shared_ptr

  if (!pipe)
    throw std::runtime_error("popen() failed!");

  std::string result;
  std::array<char, 128> buffer;
  while (!feof(pipe)) {
    if (fgets(buffer.data(), buffer.size(), pipe) != nullptr)
      result += buffer.data();
  }

  pclose(pipe);
  return result;
}

std::string sanitizeEscapes(const std::string &str) {
  std::string result;
  result.reserve(str.size());
  for (const auto c : str) {
    if (std::isprint(c))
      result.push_back(c);
  }
  return result;
}
