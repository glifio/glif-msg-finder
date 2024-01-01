glif-msg-finder
===============

A library for quickly finding messages sent to agent contracts
using the Lotus JSON-RPC API.

## CLI example

```
$ go run ./cmd --help
Find the blocks won by a miner

Usage:
  find-blocks <miner-id> [flags]

Flags:
  -h, --help             help for find-blocks
      --max-epoch uint   The minimum epoch
      --min-epoch uint   The minimum epoch
      --rpc-url string   Lotus endpoint (default "https://api.node.glif.io/rpc/v1")
      --strict           Fail if node doesn't have enough data
```

```
$ go run ./cmd f02812307
Height: 3309302
Tipset: {bafy2bzacebtoipo56i76j5olxfdulknpn5xcpdwupoanztl3vwpeypjnjw7do,bafy2bzacebufc3afhuyrnhowovcytf44uqjp4usmmtysklclp4j5a3zmdhgxu,bafy2bzacedjoncn5ycd2ooox43yeikoraqncxo3ifyckohzjtks4ha363ilsq,bafy2bzacecraa3gh5ynhfsnwluaxygtbyuh2l4exumz3aax556ws44oekafxq,bafy2bzaced3lmidtgesodtquoosl2o676bp62kjzbultee364xmpyou7rudeo,bafy2bzaceau3ii2vuye77bo7baugpjspiolck6wxfw7ptwnpzs2yceftnrj6u,bafy2bzacecrwbbd3ud7ahvmal6zjqrzctcyghenavaifqoowffffbzfz2ipde,bafy2bzaceaua5havrkcbq2cwisr72at2hupuwdpds2hsvsqbmpwkpqtjhvylk,bafy2bzaceaase5by45okiakcqpxpxwwnd7fjho57hiwg4pe6vybq6as7l54su}
Results:
3308884: 11.452573080 bafy2bzacecs3het7xj4auq6qhxcd5v6rrvwndpbc7aowzdfnltim6q6npsey6
3308815: 11.452919498 bafy2bzacebvjev42xazwwaeixsam6a6rf75jftv6ekgrfqqjb7oly6l2ogyow
3307850: 11.456476453 bafy2bzacedtywxbmzwk4tc3lqvgeh45zddjia2t3jag4dhmd65a6fbdjaacia
3307797: 11.456755700 bafy2bzacea2wq4tnqoqewd6cezklgh375t3uftmejmtdry2mf4hfm5hn4ysjw
3307592: 11.457415930 bafy2bzacedthawbiabfqxn4hjxt7yx2m6eth3fqc7os2hxktmtxyzzefblqjq
3307418: 11.457840055 bafy2bzaceax2c4xya3k3qpbsfkxhsr5toe7narrhvcbrydeaakhkx5647jzqw
Total: 68.733980716
```
