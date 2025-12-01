{
  description = "A Nix-flake-based Go 1.25 development environment";

  # Flake inputs
  inputs = {
     nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable"; # also valid: "nixpkgs"
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      goVersion = 25; # Change this to update the whole stack

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            # go (version is specified by overlay)
            go
            gotools
            golangci-lint
            cobra-cli
            mockgen
          ];

          shellHook = ''
            alias cobra="cobra-cli"
          '';
        };
      });
    };
}