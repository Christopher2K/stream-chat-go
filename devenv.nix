{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:

{
  packages = [ pkgs.git ];
  languages.go.enable = true;
  languages.go.version = "1.26.0";
}
