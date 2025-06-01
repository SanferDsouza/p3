{
  description = ''
    Practice typing passwords by comparing them to their hashes.
    Can practice multiple passwords by providing a description
    that identifies each password hash. Can be a hint.
    Never type a password incorrectly again (after some practice)!
  '';

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/24.11";
  };

  outputs =
    {
      flake-utils,
      nixpkgs,
      self,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShellNoCC {
          packages = with pkgs; [
            go
            python310
            virtualenv
          ];

          shellHook = ''
            virtualenv -p312 venv
            source ./venv/bin/activate
            pip install -qr requirements.txt
          '';
        };
      }
    );
}
