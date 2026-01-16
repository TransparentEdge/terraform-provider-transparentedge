{
  description = "Nix Flake Development Shell";

  inputs = {
    # https://www.nixhub.io/packages/go
    nixpkgs.url = "github:nixos/nixpkgs/f665af0cdb70ed27e1bd8f9fdfecaf451260fc55";
    #nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          system = system;
          config.allowUnfree = true;
        };
      in
      {
        devShell = pkgs.mkShell {
          name = "go-shell";
          buildInputs = with pkgs; [
            go
            gopls
            golangci-lint
            terraform
          ];
          shellHook = ''
            echo "shell ready"
          '';
        };
      }
    );
}
