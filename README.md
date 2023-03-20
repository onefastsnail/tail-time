# T~~ai~~l time [![CI](https://github.com/onefastsnail/tail-time/actions/workflows/ci.yml/badge.svg)](https://github.com/onefastsnail/tail-time/actions/workflows/ci.yml)

A fun wee project that creates bedtime tales using [OpenAI](https://openai.com/) and sends those tales onwards to one's reading device, for example a Kindle via email.

A very much WIP project.

I plan to make this into a usable service that would publish the generated tales for all to access, however until that time, this project can easily be forked and modified to fit your own needs. Enjoy!

## Background

I love reading to my kids. I love technology. I had an idea to mix these two together and had some fun.

## Architecture

This solution is :100: overkill, and could be done in far simpler ways, even in a single small script, but what would be the fun be in that!

![architecture.png](docs/architecture.png)
 
However, the implemented and deployed microservice architecture does provide advantages over its monolith counterpart such as:

- Increased resilience should one part of the pipeline fail.
- Easier testing, debugging and local development.
- Faster and easier future development / maintenance.
- Single and clear responsibilities of the service implementations. 
- Improved scalability should that time ever come.

All services are written in [Go](https://go.dev/).

The infrastructure is provisioned with [Terraform](https://www.terraform.io/) and deployed to [AWS](https://aws.amazon.com/) and runs fully in the [free tier](https://aws.amazon.com/free).

## Development

Coming soon.

## Up next...

Many things in the roadmap but here are a few:

- Allow the pipeline to be invoked by an Alexa custom skill.
- Get creative with the topics, maybe use OpenAI to get those.
- Experiment and explore different prompts and models. 
- Feed the tale back into OpenAI for validation and improvements.
- Build a larger tale consisting of chapters covering various topics.
- Present tales generated to the public for access and usage.
- Augment the tale with generated imaginary.
- Present the tale in prettier formats ie PDFs.
