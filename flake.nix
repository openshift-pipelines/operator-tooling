{
  description = "operator-tooling nix package";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-22.05"; # We could use nixos-unstable but.. why ?

  outputs = { self, nixpkgs }:
    let

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 self.lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs {
        inherit system;
        # Makes the config pure as well. See <nixpkgs>/top-level/impure.nix:
        config = {
          allowBroken = true;
        };
      });

    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          operator-tool = pkgs.buildGo117Module {
            pname = "operator-tool";
            inherit version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;
            subPackages = [ "cmd/operator-tool" ];

            # We use vendor, no need for vendorSha256
            vendorSha256 = null;
          };
          docker =
            let
              operator-tool = self.defaultPackage.${system};
            in
            pkgs.dockerTools.buildLayeredImage {
              name = operator-tool.pname;
              tag = operator-tool.version;
              contents = [ operator-tool ];

              config = {
                Cmd = [ "/bin/operator-tool" ];
                WorkingDir = "/";
              };
            };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.operator-tool);

      devShell = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        pkgs.mkShell {
          buildInputs = [
            pkgs.go_1_17
            pkgs.gotools
            pkgs.golangci-lint
            pkgs.gopls
            pkgs.go-outline
            pkgs.gopkgs
          ];
        });
    };
}
# }
