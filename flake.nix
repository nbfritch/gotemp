{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs }:
    let pkgs = nixpkgs.legacyPackages.x86_64-linux.pkgs;
    in
    {
      packages.x86_64-linux.default = pkgs.buildGoModule {
        pname = "gotemp";
        version = "0.1.0";
        src = ./.;
        vendorHash = null;
      };
      devShells.x86_64-linux.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          bashInteractive
          go
          gopls
        ];
      };
    };
}

