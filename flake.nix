{
  description = "Nix Flake Development Shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-22.11";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    nixpkgs-unstable,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
        pkgs-unstable = nixpkgs-unstable.legacyPackages.${system};
      in {
        devShell = pkgs.mkShell {
          name = "go-shell";
          buildInputs = with pkgs; [go gopls golangci-lint pkgs-unstable.terraform];
          shellHook = ''
            echo "shell ready"
          '';
        };
      }
    );
}
