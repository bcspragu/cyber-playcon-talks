# Cyber PlayCon 2024 talks

This repo contains the slides and associated infra for my talks. To view the talks, go to [cybr.lol](https://cybr.lol).

# Setup

## Go

[Go](https://go.dev/) is used for the live dev server. Install according to [the instructions](https://go.dev/doc/install) for your platform

## NPM

[NPM](https://www.npmjs.com/package/npm) is mainly used for the Tailwind CLI, which is total overkill, but :shrug:. NPM is usually installed via something like [nvm](https://github.com/nvm-sh/nvm).

Once you have nvm installed, run `npm i` to install dependencies

## Remark.js

Download a copy of [remark.js](https://github.com/gnab/remark), which is used for rendering Markdown + HTML slides into a web-based presentation. The server expects it to be in the repo root at `remark-latest.min.js`, you can download it from `https://remarkjs.com/downloads/remark-latest.min.js`.

# Running a presentation

To serve a presentation, run:

```bash
# Run the backend
go run . <number of talk>

# Compile CSS
npm run tailwind
```

For example, to serve talk `01-architecture-101`, you'd run:  `go run . 1`
