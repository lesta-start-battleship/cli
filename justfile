build:
  docker build -f ./build/cli.dockerfile -t "lesta-battleship-cli:dev" .

run:
  docker run -it "lesta-battleship-cli:dev"
