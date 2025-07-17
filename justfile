clientContainerName := "client"

setup-and-build:
  openssl s_client -connect battleship-lesta-start.ru:443 -showcerts </dev/null | openssl x509 -outform PEM > battleship-lesta-start.ru.crt

  docker build -f ./build/cli.dockerfile -t "lesta-battleship-cli:dev" .

@run:
  -docker rm {{clientContainerName}}
  docker run -it --name {{clientContainerName}} "lesta-battleship-cli:dev"
