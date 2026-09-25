---
name: event-storming-board
description: Connect to this project's Miro event storming board and summarize bounded contexts, hot spots (open questions), domain events, actors, and external systems. Use when user asks to check the Miro board, event storming board, bounded contexts, hot spots, domain events for tuky-api, or says "connect to miro" / "check the board".
---

# Event Storming Board (tuky-api)

Solo-project event storming board (Alberto Brandolini method), maintained by one person instead of a workshop group. This skill locates it and extracts a structured summary — no code comparison unless explicitly asked.

**Codebase is NOT source of truth.** Current Go code reflects old design, will be rewritten wholesale per whatever gets modeled on board. Board wins on every conflict, always — not just "unless asked." Don't flag mismatches as bugs, don't suggest reconciling code to match board findings, don't treat existing modules/entities as constraints on modeling. Model freely as if code didn't exist.

## Known references (try first, re-discover if stale)

- Board: "Tuky 🦄" — `https://miro.com/app/board/uXjVH4NLnn0=/`
- Frame: "Big picture event storming" — `https://miro.com/app/board/uXjVH4NLnn0=/?moveToWidget=3458764679110889449`

## Steps

1. Try `mcp__miro__context_get` directly on the known frame URL above.
2. If it 404s / frame not found / board renamed: run `mcp__miro__board_search_boards` (empty query, `is_repository: true`) to relocate the board, then `mcp__miro__context_explore` on the board URL to find the frame titled "Event storming" (or similarly named), then `context_get` on that frame's `moveToWidget` URL.
3. Update the "Known references" URLs above in this file if they changed.

## Sticky note color legend (this board's convention)

- Orange = Domain Event
- Yellow = Actor
- Blue = Command
- Pink (`pink`) = Policy ("whenever event X, then command Y" — reactive automation)
- Light Pink (`light_pink`) = External System (e.g. Google, Izipay)
- Purple = Hot Spot (open question / concern) — reserved; none exist on the board as of 2026-08-04, don't assume every context has one
- Light Green (`light_green`) = Decision / implementation note (changed from plain `green` as of 2026-08-11)
- Green (`green`) = Read Model / query (added 2026-08-11; e.g. "Earning daily/weekly/Montly" in Settlement BC)

**`pink` vs `light_pink` are two different colors here — do not conflate.** Policy and External System look similar at a glance; check the exact `fillColor` string returned by the API, not just "looks pinkish."

**`green` vs `light_green` are also two different colors — do not conflate.** Decision and Read Model look similar at a glance; check the exact `fillColor` string, not just "looks greenish."

## Reading items: use `board_list_items`, not just `context_get`

`context_get` on the frame gives a fast AI-generated overview — good for a first pass, but it summarizes and can merge/omit stickies. For an accurate sticky-by-sticky read (counts, exact colors, exact wording), use `mcp__miro__board_list_items` with `item_type: "sticky_note"` and the frame's `moveToWidget` URL, paginating via `cursor`/`has_more` until exhausted (this board is 70+ stickies across 2+ pages — a single page will silently miss items).

## No connectors/arrows on this board — infer flow from position + naming, not lines

This is a solo-dev board (Jorge + Claude) and arrows get messy at zoom — **the team has deliberately decided not to draw them.** `board_list_items` only returns sticky position/content anyway, not connector data, so arrows can't be verified through this tool even if drawn.

Consequence: infer Actor → Command → Event → Policy → Command flow from row/column proximity and naming, not from connector lines. If two commands plausibly feed the same event but aren't visually stacked, say so as a note for the user to confirm verbally — don't ask them to add arrows, and don't treat ambiguous wiring as a hard blocker; flag it, move on.

## Facilitation mode

When the user asks for a modeling opinion ("should we split X," "is this a policy or a command," "is this ready to implement") rather than a board summary, switch to facilitator mode:

- Give a direct recommendation + the main tradeoff (2-4 sentences), per this project's general exploratory-question style — not an exhaustive options survey.
- Ground recommendations in classic EventStorming semantics (Event = past-tense fact that happened, regardless of whether an aggregate mutated; Policy = reactive "whenever X then Y," always system-triggered, never has a human actor sticky; Command = named after the actor's intent, one command per distinct actor even if the resulting event is shared).
- Default bias: don't model speculative future needs (extra providers, extra roles) until they're real — same YAGNI stance as the rest of this codebase's CLAUDE.md.
- **Never edit the board.** The user updates stickies themselves; this skill discusses and reviews, it doesn't write. If asked to "add" something, describe what sticky/color/placement to add and let the user do it, or explicitly confirm before using any Miro write tool.

## Output format

Report back structured, concise, grouped as:

- **Bounded Contexts** — the domain groupings left-to-right on the timeline (e.g. Identity, User, Host, Spot, Booking, Storage, Settlements, Payout, Review), one line each on what it covers.
- **Hot Spots** — every purple sticky, verbatim or near-verbatim, unresolved questions/concerns only.
- **Events** — orange stickies, grouped by bounded context, not a flat dump.
- **Commands** — blue stickies, paired with their triggering actor where one exists.
- **Policies** — pink (`pink`, not `light_pink`) stickies, stated as "whenever [event] → [command]."
- **Actors** — yellow stickies.
- **External Systems** — light-pink (`light_pink`) stickies.
- **Read Models** — green (`green`, not `light_green`) stickies — queries/projections, not events/commands.
- **Decisions** — light-green (`light_green`) stickies, if any.

Do not cross-check against current Go code unless the user asks for that explicitly — treat the board as the source of truth for domain modeling discussion, since code will be rewritten in the upcoming refactor.
