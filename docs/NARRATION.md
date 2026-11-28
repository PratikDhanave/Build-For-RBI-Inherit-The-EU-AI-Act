# Full Narration — One Continuous Week in Bodhyanga's Banking App

Matches the live deck (`26` slides = `22` talk + `4` Q&A backups).
**Arc:** Tuesday 2:00 PM Crisis → Core Thesis → Tuesday 2:01 PM Gate → Wednesday Bug #1 → Friday Insider Bug #2 → Next Month Frankfurt Twist → Honest Gaps & Contract → Bodhyanga Close.

---

**01 · Title — 0:00 *(The 15-Second Hook)***
"Good evening. Tuesday, two o'clock in the afternoon. A customer named Bodhyanga taps one button in her banking app: *Stop sharing my data.* Behind that screen sit sixty-two AI helpers. One minute later — did they all stop? And can you prove it to a regulator in three seconds? I'm Pratik from ClearRoute. This is the story of how we built that system — and every check you see tonight, you can run right on your own laptop."

**02 · Huge Shoutout — 0:25**
"Before we jump in — a huge thank you to the FDExHyderabad operations team. A talk only looks smooth because of all the quiet work they did before any of us walked in. That's them. [small pause] Now — back to Tuesday, two o'clock."

**03 · Part 1 · Three Knocks on Our Door — 0:45 *(Tension First!)***
"Meet Bodhyanga. She taps Stop at two o'clock — and suddenly three different teams knock on our door at the exact same time. The regulators ask: *does her data stay inside India?* The risk officers ask: *she tapped Stop at two o'clock — did it actually hold at two-oh-one?* And our partner banks ask: *prove where her data went.* Three real demands. And a slide deck couldn't answer a single one."

**04 · Borders, 2:00 PM Stop & 3-Second Proof — 1:30**
"So we turned those three knocks into three everyday pictures we could actually engineer. One: a **bag tag** — Bodhyanga's Tax ID headed to a US server stops at the border. Two: a **never-erase timeline** — her Monday Yes stays locked as Seq 1, her Tuesday two o'clock Stop appends underneath as Seq 2, so the two-oh-one AI call stops cold with proof. Three: a **signed travel receipt** — answered in three seconds, not three weeks."

**05 · The Core Thesis — 2:20 *(Pause & Land This Word-for-Word)***
"[slow down, clear voice] That brings me to the one idea I want you to carry out of this room tonight: **AI compliance isn't a legal checklist — it's deterministic systems engineering. Build seventeen real controls once for India's RBI, and you inherit the EU AI Act free — because every regulator asks the exact same engineering questions.** [beat] Everything in the next thirteen minutes is runnable proof of that one sentence."

**06 · The Score, Up Front — 2:55**
"Here is the honest score before I show a single line of code: seventeen of seventeen Indian bank AI rules built in real Go. Seven European rules inherited with zero extra stack code. And two problems still open — shown openly, not hidden away. Every check you can run right on your own laptop."

**07 · 26 India Rules Split (Airport vs. Airline) — 3:30**
"When India's RBI published twenty-six AI rules in August 2025, everyone panicked. Think of an **airport versus an airline**: the airport builds the runway and control tower — that's nine national rules you cannot write app code to do. The airline engineers its own planes — that's your **seventeen controls**, across four real systems. Quote for all twenty-six and you just priced the government's job. Build seventeen."

**08 · Part 2 · Tuesday 2:01 PM: 62 Helpers Broke — One Gate Fixed It — 4:15**
"Back to Tuesday, two-oh-one PM — Bodhyanga's message hits sixty-two AI helpers. Think of **airport security**: one checkpoint at the terminal door, not an X-ray machine bolted inside all sixty-two planes. Look left — we tried the obvious way first, copying rules inside every helper. Latency blew up to two-point-four seconds, and Helper sixty-two ran last month's rules and leaked. Look right — one compiled Go Composite Policy Gate outside all sixty-two helpers checks Borders, Consent, Reason, and Role in zero-point-eight-five milliseconds."

**09 · An AI Guard Guesses. Code Rules Don't. — 5:05**
"Everyone asks: why not use another AI to judge the first AI? That's a **moody bouncer versus a metal detector**. Send Bodhyanga's exact same risky message three times: the AI judge takes up to two seconds and flips its answer — SAFE, UNSAFE, SAFE! Deterministic Go rules answer in zero-point-four milliseconds and say BLOCK — every single time."

**10 · The Whole System, on One Slide — 5:45**
"Here is the entire pipeline on one slide. Follow the green dot: Bodhyanga's cleared message passes the gate in under a millisecond, checks the permit, runs the helper, writes to the sealed diary, and emits signed reports for **both India and Europe**. Follow the red dot: any blocked call stops dead before any AI sees it and logs itself straight to the diary."

**11 · Run on Your Own Laptop — 6:25**
"Before the two bugs that nearly broke us — one fair test you can put to any AI vendor: download it and run every safety check right on your own laptop. No cloud login. No secret keys. No external API. Ours runs all rules and twenty red-team attack tests in under ten seconds."

**12 · Part 3 · Wednesday Morning (Bug #1 · Scooter License vs. Heavy Truck) — 7:05**
"Now **Wednesday morning** — bug number one. A team upgrades `var_calculator` to score Bodhyanga's credit limit. It used to be MEDIUM risk; now it touches real money — HIGH risk. Nobody revoked its old approval, and it still has sixty days left. That's a **valid scooter license, driving a heavy commercial truck**. Look at the runway: `ApprovalDriftPolicy` checks license against vehicle on every trip in under a millisecond — sees MEDIUM is not HIGH — and drops the barrier. Built for India Rule 17 — satisfies Europe Articles 9 and 14."

**13 · Live Demo 1 — 8:00 *(Frame the Human Bet First!)***
"Don't take my word for it. If I miss this one helper, an unapproved AI touches Bodhyanga's money today. Watch what happens when I press R. [press R — STOP TALKING FOR TWO FULL BEATS when `STOP var_calculator ErrRiskChanged` flashes red!] Four helpers match their permits — and `var_calculator` stops cold."

**14 · Friday Audit · Bug #2 in Twelve Words — 9:00 *(Slow Down 20%)***
"[read slowly] Now **Friday morning** — audit day. A lock that cannot catch its own team is not a lock. [HOLD TWO FULL SECONDS OF SILENCE] Protects the company. Not the customer."

**15 · Bug #2 · Three Locks (Sealed Numbered Diary) — 9:30**
"Imagine someone *inside* the bank panics before Friday's audit and quietly edits Row 2 of Bodhyanga's history — changing `advisor-7` to `advisor-99` to cover who viewed her account. Think of a **numbered diary with wax seals**: edit page two, and pages three and four tear at once. Not a blockchain — just three boring locks: SHA-256 padlock bridges between rows, an Ed25519 locked signature, and a 60-second read-only outside backup. Built for India Rule 22 — passes Europe Article 12."

**16 · Live Demo 2 — 10:30 *(Frame the Human Bet First!)***
"An insider just altered one name in Bodhyanga's history to cover their tracks. Watch whether our report lies for us — or catches us. [press R — let `BROKEN at seq 2` land] Broken at row two — and the export manifest flags `chain_verified: false`. [beat] **The report refused to lie.** [STOP TALKING FOR TWO FULL BEATS]"

**17 · Part 4 · Next Month (The Frankfurt Twist · One Aircraft Black Box) — 11:30**
"Now **next month** — the twist nobody on our team expected. Bodhyanga's family in **Frankfurt** uses the exact same banking app — and suddenly Europe's EU AI Act applies to us from Hyderabad! Leadership braced for a six-month rewrite. Instead — think of **one aircraft black box flying Hyderabad to Frankfurt**: nobody bolts on a second black box over Europe. Look at our convergence matrix: our real India RBI controls on the left flow straight into **seven EU AI Act articles** on the right with zero extra stack code — and our locked diary alone satisfies both Article 12 and Article 73."

**18 · Photograph This · 7 Safety Checks, 2 Reports — 12:30**
"Take a photo of this cheat sheet. [pause 5 seconds] Just two rows to keep: row one, our sealed diary answers India Rule 22 and Europe Article 12; and row six, telling the user 'this is AI' on the first message has been **LIVE LAW in Europe since August 2026** — and we built it eleven months earlier for India. One local command prints both signed reports."

**19 · Honest Score · 4 Built, 2 Open Gaps — 13:20**
"Now the part most conference talks skip — what we haven't solved yet. In June 2026, India's advice became strict draft rules. Four of six already worked because we engineered early. Two are honest open gaps."

**20 · Two Open Problems (Fuse vs. Plug & 10-Year Graveyard) — 13:55**
"Gap one: an **auto-fuse blows itself** — that's an accident, not a decision. The rule wants a **human hand on the plug** that stops one rogue AI mid-sentence, blocks auto-restart, and keeps Bodhyanga's payments running. Industry has solved this — **CrowdStrike Falcon Guardian** ships a four-layer stop pattern (happy to jump to Slide 23 in Q&A). Gap two is a real design clash: a **spotless live noticeboard** stays exact, while a **ten-year graveyard** of retired AIs goes stale fast. If your team solved Gap two cleanly, come find me after."

**21 · User vs. Maker — Decide Before Signing (Your Logo, Your Liability) — 14:55**
"One last rule that saves real money on Monday. Buy a vendor's AI and keep their name on it — you're a protected User. Slap your bank's logo on it — like **sticking your own label on someone else's water bottle** — or retrain it on Bodhyanga's emails, and EU Article 25 legally flips you into the Maker with full liability. Write that one line into the contract on Day 1."

**22 · Thank You — 15:45 *(Full-Circle Bodhyanga Close)***
"Three things for Monday. Build your seventeen — not all twenty-six. Put one deterministic gate outside your AIs — not sixty-two slow ones inside. And name your gaps plus that one contract line on Day 1. [beat] And **Bodhyanga**? She never saw Go, SHA-256, RBI Rule 22, or EU Article 50. She tapped Stop at two o'clock on Tuesday — and everything simply stopped. That's what deterministic systems engineering buys you: the customer never finds out how hard it was to keep the promise. Build for India. Inherit Europe. The code is open — run it on your own laptop tonight. Thank you."
[Stop. Stand still. Don't narrate the QR codes.]

---

## Q&A Backups — jump anytime with `23↵`–`26↵`

**23 · Q&A 01 · Off-Switch (`23↵`)** — "Three steps (Dimmer → Cancel Passkey → Pull the Plug) paired with CrowdStrike Falcon Guardian's four-layer stop pattern so core payments stay online."
**24 · Q&A 02 · Five Places (`24↵`)** — "One `policy.yaml`, zero country forks. India, Abu Dhabi, Dubai, Saudi, Europe checked per message."
**25 · Q&A 03 · Top 4 AI Tricks (`25↵`)** — "Prompt tricks can't sweet-talk a deterministic gate sitting outside the AI — plus 20 OWASP attack tests on every PR."
**26 · Q&A 04 · Single Point of Failure? (`26↵`)** — "No central server to crash. One rulebook in Git compiled into each of the 62 helpers — one crash isolates cleanly and fails closed while 61 keep serving."
