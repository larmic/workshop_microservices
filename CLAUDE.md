# Microservices Workshop

Educational workshop project (German language) teaching microservice architecture patterns and best practices.

## Audience & Language

- Target audience: software architects, developers becoming architects, tech leads with architectural responsibility
- All workshop content (docs, stories, instructions) is written in German
- Code, identifiers and commit messages: English

## Tech Stack

- Reference implementation in `services/` is written in Go (currently 1.25)
- Tech stack for workshop participants is intentionally **free of choice** — they may use any language/framework that fits the task
- Example domain: travel booking system (Hotel, Flight, Car rental, Booking orchestration)

## Workshop Topics

The full topic overview lives in `docs/themen.md`. Topics span foundations, resilience patterns, communication & routing, data & events, deployment & operations, and culture & organisation.

User stories for the hands-on tasks are in `docs/stories/` (`story-01` … `story-07`).

## Project Structure

- `docs/` — workshop documentation
  - `themen.md` — topic overview
  - `vorbereitung.md` — workstation setup
  - `idea.md` — workshop concept
  - `stories/` — user stories
  - `instructions/` — trainer notes
  - `questions/` — discussion prompts
- `services/` — reference Go implementation
  - `booking/` — BookingService with one folder per story (`story1` … `story7`)
  - `flight/`, `hotel/`, `car/` — domain services
  - `dashboard/` — dashboard UI
  - `traefik/` — API gateway configuration
  - `shared/` — shared libraries
  - `docker-compose*.yml`, `Makefile` — local setup
