# Pear Programming

## Inspiration

Teamwork is at the heart of all engineering disciplines, and software
engineering is no exception. It is common and productive in our industry to see
multiple engineers editing the same files simultaneously while thinking through
problems in real time. Some tools, like VS Code's Live Share, are tied to an
IDE. But like many developers, we only feel truly at home in a textual
interface. That's why we felt the need to enable real-time terminal sharing
that just works.

## What it does

Pear allows multiple people to simultaneously access a terminal as if everyone
had a keyboard plugged into the same computer, except the users can be on
opposite sides of the planet. With Pear, real-time live editing and
collaboration can extend to any tool in the terminal—which for software
engineers who use Neovim is nearly every tool in the development environment.

## How we built it

Pear is built on `libp2p`, the same networking library which powers the
[InterPlanetary File System](https://ipfs.tech/). A lightweight registry server
helps start peer-to-peer connections between clients. The entire project is
written in Go, an excellent choice for a project like this. Other libraries we
used include [cobra](https://cobra.dev/), a command-line interface builder
tool, and [bbolt](https://pkg.go.dev/go.etcd.io/bbolt), a KV database.

## Challenges we ran into

`libp2p` was not designed for ease of use. With sparse documentation,
difficult-to-discover examples, and nothing but a formal specification for
describing how its components fit together at a high level, the API was
difficult to work with. Additionally, it is not trivial for nodes to discover
each other. The most basic method, rendezvous, is not decentralized, which
defeats some of the purpose of a peer-to-peer network. The next method, mDNS,
works well, but only over local networks. Finally, there is a distributed hash
table (DHT) method, but not only is this extremely complex and somewhat
non-performant, it also does not work unless the first few nodes are pre-known
in a "bootstrap" connection. We decided to compromise by writing a simple
registry server to solve node discovery, but this proved extraordinarily
difficult itself because the process of working around the masking effected by
network address translation (NAT) is extremely complicated and equally poorly
documented.

The issue that took the longest to resolve was a subset of the node connection
process described above. What we didn't realize is that Fly.io by default gives
projects a shared IPv4 address. This doesn't work well with P2P. This was not
at all apparent given the error messages we were seeing. We solved the issue by
purchasing an upgraded version of Fly.io that gives each project its own IP
address.

## Accomplishments that we're proud of

Two of us had minimal or no Go experience, so we are proud to have learned it
so quickly. We are also proud of sticking it out and fighting through
frustrating issues with the network problems that comprised the bulk of the
challenge in this project. Overall, we engineered our way around some
significant challenges, and we're quite happy with the result.

## What we learned

Two of us learned Go. All of us learned one of the biggest and most complex
libraries we have ever seen, the Go implementation of `libp2p`. We learned
about the intricate process involved in overcoming the problems imposed by the
prevalence of NAT in modern network infrastructure. (See
https://tailscale.com/blog/how-nat-traversal-works for a good read on why this
seemingly simple problem has become very complex.)

## What's next for Pear

Without the time constraints of a hackathon, we would have liked to implement
buffers so that commands with high output do not flood the network. This should
significantly improve the app's stability. Additionally, we would have
implemented an additional layer of encryption using password-authenticated key
agreement (PAKE). Another potential feature we had considered was voice calling
capability.
