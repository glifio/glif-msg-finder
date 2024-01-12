glif-msg-finder
===============

A library for quickly finding messages sent to agent contracts
using the Lotus JSON-RPC API.

## CLI example

First, get a [Beryx Token](https://auth.zondax.ch/docs).

```
export BERYX_TOKEN=<token>

$ go run ./cmd 2
Address: 0xf0F1ceCCF78D411EeD9Ca190ca7F157140cCB2d3
Transactions:
3556643 bafy2bzacecsuexybki6qsciuon6k44qoykiics6oqwc3qu5gtpa6uvyl6wpt6 pay 0.00
3547993 bafy2bzaceamtmu4rmmno5fmn7u65fj7myu54llevwe5kh2ofooa4wcdufp7ic pay 0.00
3539343 bafy2bzacecm4uoa35sym4m6vlpj34v7g4svo3uqn3yzg5irlrkymhsl4mohyi pay 0.00
3530693 bafy2bzaceaq367ytwb3bfp5rn3bo3wgiezf2jieoy2xdv3qrz62drelmpvoga pay 0.00
3528161 bafy2bzacebpaacaw2ogsbs2yxdmyyxavnwp4uiqiefprdkfqnqw2hctxtdboo pullFunds 1.00 from f01931245
3522044 bafy2bzacecpgrq2oxhzworzhvydvb3sptckp2wwoqfhexz7uqsgug26cgm5wi pay 0.00
3513395 bafy2bzacedlsnq7azbj3c6kiu6oxnemhkf4uyziwjja3jzlmz6zsg53wpyh46 pay 0.00
3504746 bafy2bzacecghecmsimcob5qtnb7txtzesc72bilrdj7jckwipq2zn3nubu6gq pay 0.00
3496098 bafy2bzaceaqvsijdw632mk5idkkdpvdf33kmmna7i7k73rf2tvamxzfxgdh4s pay 0.00
3487449 bafy2bzacebg3nzhukr6y2apjuj4fki33mkktsewnpljsm263qn5nlt7htudaa pay 0.00
3478800 bafy2bzaceaswbr6r4besbgapxbcbakit6cwomgdjczo7foxlm4kwxuq2pxbrk pay 0.00
3470149 bafy2bzaceblxihcvraja3a3wst2qxeqengjxdrb6z6kffpkugosxekhjawgio pay 0.00
3461499 bafy2bzacedikdkyhdyyi3dynqri3toalgyfcmav2cpg3knetsvbcwzdhwn2ds pay 0.00

...

2946074 bafy2bzacebvd667zzpcs3xehlqj4hrshjn6dmxwyoxifuvi3bji3yq5jbudtc pay 0.00
2946040 bafy2bzaceb2mojywlvsx2qwyz7up6yafhftrp6txneg5zvyuaywecoboa3ob4 pay 0.00
2946027 bafy2bzacea6qaevwbdloefukwzsrzw53bsjygl4yvfrfs64a4ihrz7vw6tsly pay 0.00
2946012 bafy2bzaceaadk7rzgsmaozrfhpwvbl7kjhps5elqg7wtb6o7uz56imcir6yua pay 0.00
2945972 bafy2bzaceczjrghnd2wjc3vctvzakqqfoj4tpeuf7ust75cv42qhxrf5vs7bm pay 0.00
2945902 bafy2bzaceamsmku24ja22sxjrifdtke2slievcr2ud6kky7waakxlslq75ziq pay 0.00
2945867 bafy2bzaceaulut7kau5atd433qlpv674f2kuvi7mu6bdcy7la5tbu376rcgqa borrow 10.00
2945855 bafy2bzacedhx2k6fxro3s4xtawi7v3p256wh5s44jyf6gv3pt6ieczgsvb5o4 pay 0.00
2940946 bafy2bzaceancf3z7ccj3bn33xdkuwpetneopg5yb4bdnun477ro3o5mocuas4 pay 0.02
2895226 bafy2bzacebuouztpxe66mqqu4migtdo7rwpcr5pfhrehqwpsrxae4otn7yits confirmChangeMinerWorker f01931245
2895222 bafy2bzacedjdnvywimziojpdtabfmriwfgxe3ct2okamfw4jnmxuonfexmn6k changeMinerWorker f01931245 -> New worker: f01931242, New control addresses: [2203568]
2895073 bafy2bzaceai2osaag475dpnwvx2szvtzug7goip2wrushuug2l6v5zkbgiqbg pay 0.00
2895059 bafy2bzaceczszgayk5pamopv52ztixv25tap7cr6vpqf2eoyd5eau4v23g5wo pushFunds 2.00 to f01931245
2895051 bafy2bzacebaoqdssfqh3cte3llwilkgpdvl3me5iofrai7pyb2jhw6t3j6otg borrow 3.00
2895043 bafy2bzacec3rxcadbzdlfrzz4oykdxyidz4aeumz36oe2ztbgcoph36xhrrx2 addMiner f01931245
```

## Usage

```
$ go run ./cmd --help
Find the messages sent to an agent

Usage:
  find-messages <agent-id> [flags]

Flags:
  -h, --help             help for find-messages
      --max-epoch uint   The minimum epoch (default 18446744073709551615)
      --min-epoch uint   The minimum epoch
      --strict           Fail if node doesn't have enough data
```

## License

Apache 2

