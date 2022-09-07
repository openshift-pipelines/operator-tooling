{ pkgs }:

with pkgs;

# Configure your development environment.
#
# Documentation: https://github.com/numtide/devshell
devshell.mkShell {
  name = "operator-tooling";
  motd = ''
    Welcome to the operator-tooling application.

    If you see this message, it means your are inside the Nix shell.

    Command available:
    - nixpkgs-fmt: to format nix code
  '';
  commands = [
    {
      name = "nixpkgs-fmt";
      help = "use this to format the Nix code";
      category = "fmt";
      package = "nixpkgs-fmt";
    }
  ];

  packages = [
    go_1_17
    gotools
    golangci-lint
    gopls
    go-outline
    gopkgs
  ];
}
