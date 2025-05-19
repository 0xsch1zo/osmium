#pragma once

#include <string>

std::string exec(const std::string &cmd);

std::string sanitizeEscapes(const std::string &str);
