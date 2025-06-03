# ~/~ begin <<docs/developer-environment.md#flake.nix>>[init]
{
  description = ''
    # ~/~ begin <<docs/developer-environment.md#flake-description>>[init]
    Practice typing passwords by comparing them to their hashes.
    Can practice multiple passwords by providing a description
    that identifies each password hash. Can be a hint.
    Never type a password incorrectly again (after some practice)!
    # ~/~ end
  '';

  inputs = {
    # ~/~ begin <<docs/developer-environment.md#flake-inputs>>[init]
    flake-utils.url = "github:numtide/flake-utils";
    # ~/~ end
    # ~/~ begin <<docs/developer-environment.md#flake-inputs>>[1]
    nixpkgs.url = "github:nixos/nixpkgs/24.11";
    # ~/~ end
  };

  outputs =
    {
      # ~/~ begin <<docs/developer-environment.md#flake-output-args>>[init]
      flake-utils,
      # ~/~ end
      # ~/~ begin <<docs/developer-environment.md#flake-output-args>>[1]
      nixpkgs,
      # ~/~ end
      self,
    }:
    # ~/~ begin <<docs/developer-environment.md#flake-output-body>>[init]
    flake-utils.lib.eachDefaultSystem (
      system:
      # ~/~ begin <<docs/developer-environment.md#flake-utils-body>>[init]
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            # ~/~ begin <<docs/developer-environment.md#dev-shells-pkgs>>[init]
            go
            python312
            tmux
            virtualenv
            # ~/~ end
          ];
          shellHook = ''
            # ~/~ begin <<docs/developer-environment.md#dev-shells-shell-hook>>[init]
            virtualenv -p312 venv
            source ./venv/bin/activate
            pip install -qr requirements.txt
            # ~/~ end
          '';
        };
      }
      # ~/~ end
    );
    # ~/~ end
    
}
# ~/~ end
