# naivechain

[![wercker status](https://app.wercker.com/status/058426a6b55db2bd57c5325739cb7224/s/master "wercker status")](https://app.wercker.com/project/byKey/058426a6b55db2bd57c5325739cb7224)

a simple Blockchain inspired by https://github.com/lhartikk/naivechain

## Install

``` sh
$ go get -u github.com/m0t0k1ch1/naivechain
```

## Run

``` sh
$ naivechain [-api HTTP_SERVER_ADDRESS_FOR_API] [-p2p WEBSOCKET_SERVER_ADDRESS_FOR_P2P] [-origin P2P_ORIGIN]
```

### Example

``` sh
$ naivechain -api :3001 -p2p :6001
```

## HTTP API

### Get chain

``` sh
$ curl http://127.0.0.1:3001/blocks
```

### Mine block

``` sh
$ curl http://127.0.0.1:3001/mineBlock -d '{"data":"some data to include in the block"}'
```

### Add peer

``` sh
$ curl http://127.0.0.1:3001/addPeer -d '{"peer":"ws://127.0.0.1:6002"}'
```

### Get connected peers

``` sh
$ curl http://127.0.0.1:3001/peers
```
