# Full Narration — simple English, slide by slide

Matches the live deck (26 slides = 22 talk + 4 Q&A backups, Bodhyanga). Clock marks show when each slide should START.
Speak slowly. Every technical slide carries its real-life analogy right at the top — say the analogy first, then point at the architecture.

---

**01 · Title — 0:00**
"Good evening. I'm Pratik. This is a story about teaching a bank's AI to behave — and every check you see tonight, you can run right on your own laptop. No cloud. Nothing to hide behind. If it breaks, you will see it break."

**02 · Shoutout — 0:20**
"First, a huge thank you to the FDExHyderabad operations team. A talk only looks smooth because other people worked very hard behind the scenes. That's them. [small pause] Okay — the one idea behind everything tonight."

**03 · The Core Thesis — 0:40**
"[read slowly, clearly] AI compliance isn't a legal checklist — it's deterministic systems engineering. Build seventeen real controls once for India's RBI, and you inherit the EU AI Act free — because every regulator asks the exact same engineering questions. [beat] Every diagram and live Go demo in the next fifteen minutes is proof of that one sentence."

**04 · The Ending, Up Front — 1:15**
"Here is the honest score before I show a single line of code: seventeen of seventeen Indian bank AI rules built in real code. Seven European rules inherited with zero extra stack code. And two problems still open — shown openly, not hidden away. Every check you can run right on your own laptop."

**05 · Part 1 · Three Knocks — 1:50**
"Meet Bodhyanga. She uses a banking app. Behind that app sit sixty-two AI helpers. One Tuesday afternoon, three teams knocked on our door. Regulators asked: does her data stay in the country? Risk officers asked: she tapped Stop at two o'clock — did it actually stop? Partner banks asked: prove where her data went. A slide deck couldn't answer a single one."

**06 · Borders, 2:00 PM Stop & 3-Second Proof — 2:30**
"Three everyday pictures that become three engineering controls. One: a bag tag — the wrong country tag never boards. Two: a never-erase timeline — her Monday Yes stays locked as Seq 1, her Tuesday two o'clock Stop appends below as Seq 2, so a two-oh-one AI call stops cold. Three: a signed travel receipt — answered in three seconds, not three weeks."

**07 · 26 India Rules Split (Airport vs. Airline) — 3:20**
"Then India's central bank published twenty-six AI rules. Think of an airport versus an airline: the airport builds the runway and control tower — that's nine national rules you cannot code in an app. The airline engineers its own planes — that's your seventeen controls, built into four real systems. Quote for all twenty-six and you just priced the government's job. Build seventeen."

**08 · Part 2 · 62 Helpers Broke — One Gate Fixed It (Airport Security) — 4:10**
"How do you police sixty-two AI helpers? Think of airport security: one checkpoint at the terminal door, not an X-ray machine bolted inside all sixty-two planes. Look left — copying rules inside every helper turned into a tangled mesh hitting two-point-four seconds, and Helper sixty-two ran last month's rules and leaked. Look right — one compiled Go Composite Policy Gate outside all sixty-two helpers checks Borders, Consent, Reason, and Role in zero-point-eight-five milliseconds."

**09 · An AI Guard Guesses. Code Rules Don't. — 5:05**
"People ask: why not use another AI to judge the first AI? That's a moody bouncer versus a metal detector. Watch the exact same risky message sent three times: the AI judge takes up to two seconds and flips its answer — SAFE, UNSAFE, SAFE! Deterministic Go rules answer in zero-point-four milliseconds and say BLOCK — every single time."

**10 · The Whole System, on One Slide — 5:50**
"Here is the entire pipeline on one slide. Follow the green dot: Bodhyanga's message passes the gate, checks the AI's permit, runs the helper, writes to the sealed diary, and outputs signed reports for both India and Europe. Follow the red dot: any blocked call stops dead before any AI sees it and logs itself straight to the diary."

**11 · Run on Your Own Laptop — 6:30**
"One simple test you can use on any AI vendor: download it and run every safety check right on your own laptop. No cloud login. No secret keys. No external API. Ours runs all rules and twenty red-team attack tests in under ten seconds — you can clone it and verify tonight."

**12 · Part 3 · Bug #1 · Old Permit (Scooter License vs. Heavy Truck) — 7:15**
"Part three — two production bugs that caught us off guard. Bug one: imagine passing your driving test on a scooter. Your license is valid, never expired — then one day you drive a heavy commercial truck with that scooter license in your pocket. Look at the runway: var_calculator was approved as MEDIUM risk with sixty days left, then upgraded to HIGH risk touching real money. ApprovalDriftPolicy checks license against vehicle on every trip in under a millisecond — sees MEDIUM is not HIGH — and drops the barrier. Valid license. Wrong vehicle. Built for India Rule 17 — satisfies Europe Articles 9 and 14."

**13 · Live Demo 1 — 8:15**
"Don't take my word for it — real Go code running right here in the browser, and you can run it on your own laptop tonight. [press R] Five AI helpers ask to run. Watch: payments clears, expired and rejected stop — and var_calculator stops cold on ErrRiskChanged. [stop talking — two full beats]"

**14 · Bug #2 in Twelve Words — 9:15**
"[read slowly] A lock that cannot catch its own team is not a lock. [silence, two seconds] Protects the company. Not the customer."

**15 · Bug #2 · Three Locks (Sealed Numbered Diary) — 9:45**
"Think of a numbered paper diary where every page copies the wax seal of the page before it. Quietly edit page two — pages three and four tear at once. Look at the code: not a blockchain, just three boring locks — SHA-256 linked padlocks between rows, an Ed25519 locked signature, and a 60-second read-only outside backup. Watch the cursor edit advisor-7 to advisor-99 on Seq 2: the padlock snaps open and the chain breaks immediately. Built for India Rule 22 — passes Europe Article 12."

**16 · Live Demo 2 — 10:45**
"Again — real Go code right here. [press R] Three rows write cleanly. Then the insider edits row two: advisor-7 becomes advisor-99. Watch the verify check: BROKEN at seq 2 — and the export manifest refuses to hide it. [beat] The report refused to lie. [stop talking — two full beats]"

**17 · Part 4 · Build 17 for India — Get 7 for Europe Free (One Black Box) — 11:45**
"Now the payoff of our opening thesis. When a plane flies Hyderabad to Frankfurt, nobody bolts on a second black box over Europe. One engineered safety core passes inspection in both places. Look at the convergence matrix: real India RBI controls on the left flow straight into seven EU AI Act articles on the right with zero extra stack code — and our locked diary alone satisfies both Article 12 and Article 73. Build once. Inspected twice."

**18 · Photograph This · 7 Safety Checks, 2 Reports — 12:45**
"Take a photo of this cheat sheet. [pause] Two rows to remember: row one, our locked diary answers both India Rule 22 and Europe Article 12; and row six, telling the customer 'this is AI' on the first message has been LIVE LAW in Europe since August 2026 — and we built it eleven months earlier for India. One local command prints both signed reports."

**19 · Honest Score · 4 Built, 2 Open Gaps — 13:40**
"Now the part most talks skip — what we haven't solved yet. In June 2026, India's advice became strict draft rules. Four of six already worked because we engineered early. Two are honest open gaps."

**20 · Two Open Problems (Fuse vs. Plug & 10-Year Graveyard) — 14:20**
"Gap one: an auto-fuse blows itself — that's an accident, not a decision. The rule wants a human hand on the plug to stop one AI mid-sentence, block auto-restart, and keep Bodhyanga's payments running. Industry has solved this — CrowdStrike Falcon Guardian ships a four-layer stop pattern (I have the backup diagram on Slide 23 if you want to dig in). Gap two is a real design tension: a spotless live noticeboard stays exact, while a ten-year graveyard of retired AIs goes stale fast. If your team solved Gap two cleanly, find me after."

**21 · User vs. Maker — Decide Before Signing (Your Logo, Your Liability) — 15:20**
"One last rule that saves your customer real pain on Monday. Buy a vendor's AI and keep their name on it — you're a protected User. Slap your own logo on it — like sticking your own label on someone else's water bottle — or retrain it on customer data, and EU Article 25 flips you into the Maker with full legal liability. Put that one line in the contract on Day 1."

**22 · Thank You — 16:15**
"Three things for Monday. Build your seventeen — not all twenty-six. Put one deterministic gate outside your AIs — not sixty-two slow ones inside. And name your gaps plus that one contract line on Day 1. [beat] And Bodhyanga? She never saw any of it. Her data stayed in India. Her two o'clock Stop held at two-oh-one. The bank got signed proof in three seconds. That's what deterministic systems engineering buys you: the customer never finds out how hard it was to keep the promise. Build for India. Inherit Europe. The code is open — run it on your own laptop tonight. Thank you."
[Stop. Stand still. Don't narrate the QR codes.]

---

## Q&A backups — jump anytime with `23↵`–`26↵`

**23 · Q&A 01 · Off-Switch (`23↵`)** — "Great question — this is open gap number one. Three steps (Dimmer → Cancel Passkey → Pull the Plug) paired with CrowdStrike Falcon Guardian's four-layer stop pattern so core payments stay online."

**24 · Q&A 02 · Five Places (`24↵`)** — "One `policy.yaml`, one app, zero country copies. India, Abu Dhabi, Dubai, Saudi, Europe — each jurisdiction is declarative config checked per message."

**25 · Q&A 03 · Top 4 AI Tricks (`25↵`)** — "Prompt tricks can't sweet-talk the gate because the gate sits outside the AI's brain — and twenty OWASP-tagged attack tests run automatically on every PR."

**26 · Q&A 04 · Single Point of Failure? (`26↵`)** — "No central server to crash. One rulebook in Git compiled into each of the sixty-two helpers — one crash isolates cleanly and fails closed while sixty-one keep serving."
