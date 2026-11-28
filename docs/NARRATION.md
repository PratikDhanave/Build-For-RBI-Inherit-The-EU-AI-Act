# Full Narration — simple English, slide by slide

Matches the live deck (32 slides, Bodhyanga). Clock marks show when each slide should START.
Speak slowly. The full stops are part of the script.

---

**01 · Title — 0:00**
"Good evening. I'm Pratik. This is a story about teaching a bank's AI to behave — and every check you see tonight, you can run right on your own laptop. No cloud. Nothing to hide behind. If it breaks, you will see it break."

**02 · Shoutout — 0:20**
"First, a huge thank you to the FDExHyderabad operations team. A talk only looks smooth because other people worked very hard behind the scenes. That's them. [small pause] Okay — the story."

**03 · The ending, up front — 0:45**
"I'll tell you the ending first, so you know where we're going. We built all seventeen of India's bank AI safety rules — in real code. Seven of Europe's rules passed automatically — zero extra work. And two problems are still open. I'll show you those too, openly. Every check you can run right on your own laptop. You should not have to take a speaker's word for anything."

**04 · Part 1 · Three knocks — 1:20**
"Meet Bodhyanga. She uses a banking app. Behind that app sit sixty-two AI helpers. One Tuesday afternoon, three different teams knocked on our door. The regulators asked: does her data stay in the country? The risk team asked: she tapped Stop at two o'clock — did it actually stop? And the partner banks asked: prove where her data went. Three questions. And a slide deck couldn't answer a single one."

**05 · Borders & the 2:00 PM stop — 1:55**
"So here are three everyday pictures. One: a bag tag at the airport — the wrong tag never flies. Two: a notebook written in pen — we never erase her Yes; we write her Stop underneath it. New line. Never an eraser. Three: a signed travel receipt — proof of where her data went, in three seconds, not three weeks. Don't memorize these. Each one comes back when we build it."

**06 · 26 India rules — 2:30**
"Then India's central bank published twenty-six AI rules. Before you build anything — think of an airline. An airline never builds the runway or the control tower. That's the airport's job — nine of these rules. The airline looks after its own planes — that's your seventeen. Seventeen rules, four systems: a permit check, a sealed diary, a live list, and a simple notice that says 'this is AI.'"

**07 · Airport vs Airline — 3:10**
"That's the whole idea in one picture. The airport builds the city's infrastructure. You keep your own planes safe. [point] Permit check. Black box. Live fleet list. Cabin notice. If you quote for all twenty-six, you just priced the government's own job. Quote for seventeen."

**08 · Part 2 · One gate — 3:50**
"Now — how do you check sixty-two AI helpers? Think of airport security. One checkpoint at the door — not an X-ray machine bolted inside every plane. We tried the obvious way first: rules inside every helper. It got slow — two and a half seconds. Copies went stale. Helper sixty-two ran last month's rules and leaked. So we pulled everything out into one gate. Every message passes it in under a millisecond. Sixty-two planes. One door."

**09 · Airport Security — 4:30**
"That's all this slide is. One checkpoint at the entrance. The planes stay simple inside. When a new rule arrives next month, you change it in one place — not sixty-two."

**10 · Rules vs AI Judge — 5:05**
"People ask: why not use another AI to judge the first AI? We tried. Watch — the same risky message, sent three times. The AI judge takes up to two seconds, and changes its mind: safe, unsafe, safe. The same message! Plain code answers in under a millisecond and says BLOCK — every single time. An AI guard guesses. Code repeats."

**11 · Whole system, one slide — 5:45**
"Here's everything together. Follow the green dot: her message passes the gate, the permit is checked, the helper answers, the diary writes a line, and signed reports come out — one for India, one for Europe. Follow the red dot: a bad message stops dead at the gate, and even the stop is written in the diary."

**12 · Run on your laptop — 6:20**
"Before the two bugs — one simple test you can use on any AI vendor. Say: download it and run every safety check right on your own laptop. No cloud login. No secret keys. Ours finishes in under ten seconds — every check you've seen tonight, you can run yourself in five minutes. A safety check nobody can run is just a claim."

**13 · Part 3 · Scooter vs Truck — 7:00**
"Part three. Two bugs that caught us off guard. First — a story. You pass your driving test on a scooter. Your license is valid. Never expired. Never cancelled. Then one day you climb into a huge truck — with that scooter license in your pocket. [beat] Watch the barrier slam down. Keep that truck in your head —" [click]

**14 · Bug #1 · Old Permit — 7:40**
"— here it is, in our app. The license says scooter: approved as MEDIUM risk, sixty days left. But the vehicle today is a truck: the team upgraded this AI to HIGH risk. The gate checks the license against the vehicle on every single trip — in under a millisecond — sees MEDIUM is not HIGH, and drops the barrier. Valid license. Wrong vehicle."

**15 · Live Demo 1 — 8:15**
"Don't take my word for it. This is real code, running right here — and you can run it on your own laptop tonight. [press R] Five AI helpers ask to run. Watch. [narrate the result] Four pass — their permits match. And var_calculator — blocked. The permit says MEDIUM, the live helper is HIGH. The gate caught it. [stop talking — two full beats]"

**16 · Bug #2 in twelve words — 9:20**
"[read slowly] A record book your own team can quietly edit protects the company — not the customer. [silence, two seconds] Protects the company. Not the customer."

**17 · Sealed Diary Pages — 9:50**
"Second story. A numbered paper diary. Every page copies the wax seal of the page before it. Someone inside quietly edits page two. Pages three and four tear — at once. And if they rewrite the whole book from page one? It clashes with the copy we keep outside, made every sixty seconds. Now the same diary, in code —" [click]

**18 · Three locks — 10:30**
"— not a blockchain. Just three simple locks. Every row copies the seal of the row before it. A locked signature nobody can fake. And a read-only copy saved outside every sixty seconds. Watch the cursor change advisor-7 to advisor-99 — and everything after it tears, just like paper. Edit yesterday — and today tears."

**19 · Live Demo 2 — 11:05**
"Again — real code, right here. [press R] The diary writes three lines. Now the attack: we quietly edit line two, advisor-7 becomes advisor-99. Run the check. [narrate] Broken — at line two. The diary caught its own team. [beat] The report refused to lie. [stop talking — two full beats]"

**20 · Part 4 · One Black Box — 12:15**
"Part four. The European surprise. When a plane flies Hyderabad to Frankfurt, nobody bolts on a second black box over Europe. One set of safety checks passes inspection in both places. One black box, two inspectors —" [click]

**21 · 17 India → 7 Europe — 12:55**
"— and here's why Europe matters even if you sit in Hyderabad. Europe's AI law applies the moment anyone in Europe uses what your AI produces. We braced for a second huge project. Instead — seven of our seventeen Indian checks passed Europe's law with zero extra code. Same four checks: stamped in India on the left, passing Europe on the right. Build once. Inspected twice."

**22 · Photograph this — 13:30**
"Take a photo of this slide. [pause — let them] Two rows only. Row one: our sealed diary satisfies India Rule 22 and Europe Article 12 — one diary, two laws. And the bottom row: telling the user 'this is AI' on the first chat has been live law in Europe since August 2026. We built it for India eleven months earlier. [pause for photos]"

**23 · Honest score — 14:25**
"Now the part most conference talks skip — what we got wrong. In June 2026, India's advice became strict draft rules. We graded our own code honestly: four already pass, because we started early. Two are real, open gaps. We're not hiding them — they're the next slide."

**24 · Pull the Plug & the Graveyard — 15:05**
"Gap one, as a story: a fuse blows by itself. That's an accident — not a decision. The rule wants a human hand on the plug: dim it, cancel its badge, pull the plug — while the lights next door stay on. Bodhyanga's payments keep running. Gap two: a spotless live noticeboard — versus a ten-year graveyard of retired AIs. A fuse is an accident. A plug is a decision."

**25 · Two hard problems — 15:45**
"Why aren't these weekend fixes? The off-switch must stop an AI mid-sentence, stop it restarting itself, and keep payments online — all three at once. And the graveyard is a real design clash: a live list never goes stale, but ten years of stored records always do. These are genuinely hard. If your team has solved either one cleanly — please find me after this talk. I mean it."

**26 · The Legal Tripwire — 16:40**
"One last warning — it can save your customer a very expensive mistake. You buy an AI tool from a vendor. Keep the vendor's name, adjust the settings — you're just the user. Safe. But put your own logo on it, or retrain it on your customers' emails — and you just stepped over an invisible legal tripwire. You became the maker. All the maker's legal duties — now yours. Your logo. Your liability."

**27 · Put it in the contract — 17:30**
"So write one line into the contract, on day one: 'If we rename it or retrain it, we take on the maker's duties — and we price for that.' One line. Signed before any work starts. That line has saved real money and real pain. And one small footnote for banks: even as just a user, credit and insurance AI still needs a fundamental-rights check in Europe."

**28 · Thank you — 18:00**
"Three things for Monday. Build your seventeen — not all twenty-six. Put one fast gate outside your AIs — not sixty-two slow ones inside. And name your gaps before a regulator names them for you — then put that one line in the contract. [beat] And Bodhyanga? She never saw any of it. Her data stayed home. Her two o'clock Stop held at two-oh-one. The bank got its proof in three seconds. She tapped one button — and everything simply worked. That's the whole job: the customer should never find out how hard it was to keep the promise. We built it for India. Europe came free. The code is open — it's yours tonight. Thank you."
[Stop. Stand still. Let them clap. Don't narrate the QR codes.]

---

## Q&A backups — opening lines (only if asked)

**29 · Off-Switch** — "Great question — this is open gap number one. Here's the honest ladder: slow it down, cancel its badge, pull the plug. Industry does the same — CrowdStrike ships a four-layer stop. The hard part is keeping payments alive while you pull."

**30 · Five Places** — "One rules file, one app, zero country copies. India, Abu Dhabi, Dubai, Saudi, Europe — each country is a few lines of configuration, checked per message. Expanding to a new country never forks the codebase."

**31 · Top 4 AI Tricks** — "Clever prompts can't sweet-talk the gate, because the gate sits outside the AI's brain. And twenty attack tests run automatically on every code change — the build fails if any slips through."

**32 · Single point of failure?** — "Fair challenge — but there's no central server. One rulebook in Git, compiled into each of the sixty-two helpers. One helper crashes, sixty-one keep serving. And if the rules file is ever broken, the gate fails closed — it blocks, it never shrugs."
