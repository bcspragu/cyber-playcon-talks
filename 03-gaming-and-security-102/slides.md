class: center, middle, intro

<h1 class="text-7xl mb-0">Gaming & Security 102</h1>

.text-3xl[Brandon Sprague]

2024-08-02

???

Gist: Talking about all the ways you can be hacked while playing games
Specifically...

---

class: contain
background-image: url(/assets/minecraft.jpg)

<h1 class="text-center -mt-1">Minecraft</h1>

---

class: center, middle

# Different types of hacks

???

Not all attacks are the same, I bucket them into a few different categories.

---

# Social Engineering attacks


Basically just people trying to trick you.

<div class="text-center text-4xl my-20">
"Hello I am from PlanetMinecraft, you have a wirus, please give me admin"
</div>

<a class="text-right block" href="https://github.com/wodxgod/Griefing-Methods/tree/master/Social%20Engineering">Some more sophisticated examples</a>

???

- Can involve people trying to get you to talk to them outside of Minecraft
- Usually has some urgency to their messaging
- The attack is mostly about getting you to give them something, like an account password, or downloading something, like a plugin or program.
- Also talk about credential stuffing

---

# Preventing Social Engineering attacks

- Don't download random things (including plugins) from the internet
- Use strong, unique passwords and a password manager
- Don't take candy from strangers
- For games with many servers, don't play on sketchy servers

---

# "Skript kiddie" attacks

People downloading malicious tools off the internet and running them.

Some examples:

- Renting Botnets + DDoS attacks
  - Sold as "stress tester" tools online
- Or protocol-level stuff: SYN-floods, Ping of Death, Slow Loris
- Generic tools you point at a server
  - To look for vulns/attack surface: [Spiderfoot](https://github.com/smicallef/spiderfoot), [Nuclei](https://github.com/projectdiscovery/nuclei), [gobuster](https://github.com/OJ/gobuster)
  - To control hacked servers: [Havoc](https://havocframework.com/)
  - A little bit of everything: [Metasploit](https://www.metasploit.com/)

---

# Preventing Skript kiddie attacks

.text-center.block[Try not to make enemies with people who are likely to be "skript kiddies".]

If that's unavoidable:

- Keep your software up-to-date
- Don't expose anything to the internet that doesn't need to be online
- Don't click sketchy links or download random things

---

# Actual, novel exploits

This is when someone finds a new weakness in a piece of software, usually called a "zero-day" because it hasn't been patched yet.

.text-center.block[This is what I think of when talking about "hacking".]

Thse weaknesses can take all sorts of different shapes:

- **[Disruptive]** Taking down a service
- **[Problematic]** Extracting sensitive data
- **[Uh oh]** Remote Code Execution (RCE)

???

That sounds pretty bad, how do you defend against it?

---

class: center, middle

# Defending against novel exploits

---

class: center, middle

...you don't.

???

- By their nature, these are *new* attacks.
- There may be mitigations against them, depends on the attack.

---

# My Favorite Example

.text-center.block.text-5xl[[RANDAR](https://github.com/spawnmason/randar-explanation/tree/master)]

Uses the [Lenstra–Lenstra–Lovász (LLL) lattice basis reduction algorithm](https://en.wikipedia.org/wiki/Lenstra%E2%80%93Lenstra%E2%80%93Lov%C3%A1sz_lattice_basis_reduction_algorithm) to find users locations in a Minecraft server.

<img class="w-[70%] mx-auto" src="/assets/lll.png" /img>

???

Basically, someone realized that Minecraft was using the same random number generator for two unrelated things: placing items in the world and item drops.

---

# My Favorite Example

<img src="/assets/randar.svg" /img>

---

# Another Cool Hack, but in Factorio

<img src="/assets/factorio.png" /img>
<a class="text-center block" href="https://memorycorruption.net/posts/rce-lua-factorio/">Source</a>

???

- Previous hack allowed hacker to see where other players were on the map
- This hack allows a malicious server to execute arbitrary code on any player's machine

---

# A Less Fun Example: Log4Shell

[Log4Shell](https://en.wikipedia.org/wiki/Log4Shell) was one of the largest vulnerabilities in recent memory, affected a large portion of software written in Java.

<img src="/assets/log4j-minecraft.png" /img>
<a class="text-center block" href="https://www.pcmag.com/opinions/critical-exploit-for-apache-log4j2-could-be-far-reaching-proves-real-in">Source</a>

???

- Opposite of the Factorio example, which was server compromising client
- This is a player compromising a Minecraft server with a well-crafted chat message.

---

# Brief Tangent: Gaming Architectures

<img src="/assets/client-server-p2p.jpg" /img>
<a class="text-xs text-center block" href="https://community.fs.com/article/client-server-vs-peer-to-peer-networks.html">(Image Source)</a>

???

- Minecraft is a client-server game, which means all players communicate with a server
- In P2P games, all players' computers communicate directly with each other.
  - From a security perspective, this means your IP might be exposed
  - Not the worst thing in the world, but worth noting

---

class: center, middle

# What have we learned?

???

- Main takeaways around staying safe while gaming?

---

class: center, middle

# Questions?
