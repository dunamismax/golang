// Copyright (c) 2025-present dunamismax. All rights reserved.
//
// filename: main.go
// author:   dunamismax
// version:  1.0.1
// date:     06-19-2025
// github:   <https://github.com/dunamismax>
// description: A single-file, CLI-based text adventure game: Veridian Echo.
//

package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// --- Data Models ---

// Player represents the state of the player character, Kael.
// It holds the current location, inventory for tracking key items or story flags,
// and the chosen narrative path.
type Player struct {
	CurrentLocationKey string
	Inventory          map[string]bool // Tracks key items, data, or flags.
	Path               string          // Vengeance, Liberation, Self-Preservation
}

// Choice represents a single decision the player can make at a given Location.
// It contains the descriptive text for the choice, the key for the next
// Location, and an optional Action function to modify the player's state.
type Choice struct {
	Text            string
	NextLocationKey string
	Action          func(p *Player)
}

// Location represents a single scene, area, or state within the game.
// It contains the descriptive text for the scene and the available Choices
// the player can make.
type Location struct {
	Title       string
	Description string
	Choices     []Choice
}

// Game encapsulates the entire state of the application. It holds references
// to the player's state, the world's locations, the main loop's running
// status, and the structured logger.
type Game struct {
	Player    *Player
	Locations map[string]*Location
	IsRunning bool
	Logger    *slog.Logger
}

// --- Story Content ---

// initializeGame creates the game world, defining all locations and narrative arcs.
// It acts as the single source of truth for the game's content.
func initializeGame(logger *slog.Logger) *Game {
	player := &Player{
		CurrentLocationKey: "start",
		Inventory:          make(map[string]bool),
	}

	locations := map[string]*Location{
		// --- ACT I: The Inciting Incident ---
		"start": {
			Title: "[SYSTEM_INITIALIZING]",
			Description: `[> DECRYPTION_KEY_ACCEPTED]
[> BOOTING_NEURAL_INTERFACE]
Welcome, Shadowrunner.

The rain never stops in Veridia. It slicks the streets with a neon sheen, reflecting the monolithic OmniCorp tower that pierces the smog-choked sky. You are Kael, a Static, an outcast. Once a promising OmniCorp security programmer, now just another ghost in the machine's shadow. Tonight, you're on a routine data heist, targeting a low-level corporate server. Easy credits. Or so you thought.`,
			Choices: []Choice{
				{Text: "Breach the server.", NextLocationKey: "heist_success"},
			},
		},
		"heist_success": {
			Title:       "The Heist",
			Description: `Your fingers dance across the holographic interface of your deck. The server's defenses are laughably weak. You bypass the outer layers of ICE (Intrusion Countermeasures Electronics) with practiced ease. Inside, amongst the mundane corporate memos and shipping manifests, something feels... different. A single data fragment, cloaked in an encryption you've never seen before. It's alien, impossibly complex.`,
			Choices: []Choice{
				{Text: "Download the anomalous fragment.", NextLocationKey: "fragment_acquired"},
				{Text: "Ignore it and take the standard paydata.", NextLocationKey: "game_over_static"},
			},
		},
		"fragment_acquired": {
			Title:       "Return to the Shadows",
			Description: `Back in the relative safety of your cramped apartment in the Neon Slums, you slot the data fragment into your deck. Before you begin the arduous task of decrypting it, you take a moment to collect your thoughts. The rain patters against your window, a familiar, dreary rhythm.`,
			Choices: []Choice{
				{Text: "Settle in and decide what to do next.", NextLocationKey: "apartment_hub"},
			},
		},
		"apartment_hub": {
			Title:       "Kael's Apartment",
			Description: `Your sanctuary. A single room cluttered with old tech, data cables, and the ghosts of your past. Your deck hums on the workbench, the anomalous fragment waiting patiently within its core memory.`,
			Choices: []Choice{
				{Text: "Work on decrypting the anomalous data fragment.", NextLocationKey: "lyra_deal"},
				{Text: "Check the public news feeds.", NextLocationKey: "apartment_news_feed"},
				{Text: "Look out the window.", NextLocationKey: "apartment_window_view"},
				{Text: "Check your personal messages.", NextLocationKey: "apartment_messages"},
			},
		},
		"apartment_news_feed": {
			Title:       "OmniCorp News Network",
			Description: `You tune into the official ONN feed. A smiling anchor talks about another record quarter for OmniCorp's atmospheric purification program. "Cleaner air, brighter futures!" she chirps, while sponsored footage of gleaming towers scrolls by. A ticker at the bottom mentions 'minor Static-related disruptions in the lower sectors successfully contained by OmniCorp Enforcers'.`,
			Choices: []Choice{
				{Text: "Return to the main room.", NextLocationKey: "apartment_hub"},
			},
		},
		"apartment_window_view": {
			Title:       "Window View",
			Description: `You gaze out at the perpetual twilight of the Neon Slums. Below, rivers of neon light from storefronts and noodle stands reflect off the wet pavement. Figures huddle under awnings, their faces obscured by shadows and steam. In the distance, the OmniCorp tower looms, a black spike against a bruised purple sky, its red signal light blinking like a malevolent eye.`,
			Choices: []Choice{
				{Text: "Step away from the window.", NextLocationKey: "apartment_hub"},
			},
		},
		"apartment_messages": {
			Title:       "Personal Messages",
			Description: `You open your inbox. Besides the usual junk mail for black market cybernetics, there's a single, encrypted message from 'Rizzo', an old contact from your past life. It's a week old. It reads: "Heard they pinned the breach on you. Keep your head down, Kael. They're not just kicking people out anymore. They're making them disappear."`,
			Choices: []Choice{
				{Text: "Close the inbox.", NextLocationKey: "apartment_hub"},
			},
		},
		"lyra_deal": {
			Title:       "The Deal",
			Description: `Lyra's voice is a strange mix of synthetic tones and genuine emotion. "OmniCorp is building something, Kael. A final solution for dissent they call 'Project Chimera'. It's the real reason you were cast out. They're turning the Nexus from a tool of control into a weapon of psychic annihilation. Help me expose them, and I'll give you everything you need to clear your name and watch the ones who wronged you burn."`,
			Choices: []Choice{
				{Text: "[Path of Vengeance] Agree. You'll make them pay.", NextLocationKey: "vengeance_intro"},
				{Text: "[Path of Liberation] Agree. This is bigger than you. The city needs to be free.", NextLocationKey: "liberation_intro"},
				{Text: "[Path of Self-Preservation] Refuse. This is too dangerous. Sell the data.", NextLocationKey: "preservation_intro"},
			},
		},

		// --- PATH OF VENGEANCE ---
		"vengeance_intro": {
			Title:       "Path of Vengeance: First Target",
			Description: `"Good," Lyra says, a cold satisfaction in her tone. "Your former supervisor, Marcus Thorne. He delivered the kill order on your career. He's at the 'Synapse' club tonight, in the VIP lounge. His terminal is his life. Get access, and we'll own him." You feel a grim smile touch your lips. Revenge is a powerful motivator.`,
			Choices: []Choice{
				{Text: "Infiltrate the 'Synapse' club.", NextLocationKey: "vengeance_club"},
			},
		},
		"vengeance_club": {
			Title:       "The Synapse Club",
			Description: `The Synapse club is a sensory overload of chrome, bass, and illegal cybernetics. Thorne is exactly where Lyra said he'd be, engrossed in a high-stakes data-poker game. His personal terminal sits on the table beside him. You need a way to get him to leave it unguarded.`,
			Choices: []Choice{
				{Text: "Bribe a server to spill a drink on him.", NextLocationKey: "vengeance_distraction_success"},
				{Text: "Attempt to social engineer him directly.", NextLocationKey: "vengeance_social_fail"},
			},
		},
		"vengeance_distraction_success": {
			Title: "Hacking Minigame: Data Skim",
			Description: `The server 'accidentally' spills a tray of vibrant blue liquid on Thorne. He curses and storms off to the washroom, leaving his terminal. You have 60 seconds. You slide into his seat and jack in.
[HACKING CHALLENGE]
Lyra: "His files are secured with biometric-linked encryption. I can't break it, but I can clone the security token while it's active. Run the 'clone_sec_token.sh' script. Don't get fancy!"`,
			Choices: []Choice{
				{Text: "run clone_sec_token.sh", NextLocationKey: "vengeance_chapter_end"},
				{Text: "run brute_force_decrypt.exe", NextLocationKey: "game_over_heist_caught"},
			},
		},
		"vengeance_chapter_end": {
			Title:       "Vengeance: The First Blow",
			Description: `The token clones successfully. You jack out just as Thorne returns. He doesn't suspect a thing. Back in your apartment, you and Lyra now have the key to his digital life. The first piece of your revenge is in place.`,
			Choices: []Choice{
				{Text: "Use the cloned token to breach Thorne's drive.", NextLocationKey: "vengeance_chapter_2_start"},
			},
		},
		"vengeance_chapter_2_start": {
			Title:       "Chapter 2: The Digital Vault",
			Description: `Back in the neon glow of your apartment, the cloned token pulses on your deck. "The token is live," Lyra states. "It's a key into Thorne's personal OmniCorp network drive. He won't even know we're there... if we're careful. I'm routing our connection through a series of ghost nodes. Let's see what secrets he's hiding."`,
			Choices: []Choice{
				{Text: "Initiate the connection.", NextLocationKey: "vengeance_thorne_drive"},
			},
		},
		"vengeance_thorne_drive": {
			Title: "Thorne's Private Drive",
			Description: `The connection solidifies. You're in. Thorne's file structure is pristine, organized. Folders for 'Financials', 'Personnel Reviews', and 'Family Photos'. But one folder stands out, its name sending a chill through you: 'Project_Chimera_Dossier'.
"That's it, Kael," Lyra whispers. "That's the project I told you about. Whatever's in there is the key."`,
			Choices: []Choice{
				{Text: "Access the 'Project_Chimera_Dossier' folder.", NextLocationKey: "vengeance_jax_intrusion"},
			},
		},
		"vengeance_jax_intrusion": {
			Title: "Intrusion Detected",
			Description: `As you're about to open the folder, a new process flashes in the system monitor: 'user_JAX_sec.analyst' has just logged into the same drive.
"Wait!" Lyra warns, her voice sharp. "Someone else is here! An OmniCorp security analyst. He's running a trace! We're out of time. You need to grab the data and scrub our entry log before he finds us!"`,
			Choices: []Choice{
				{Text: "Prepare for a rapid exfiltration.", NextLocationKey: "vengeance_jax_puzzle"},
			},
		},
		"vengeance_jax_puzzle": {
			Title: "Hacking Minigame: Race Against the Trace",
			Description: `A timer appears on your interface, counting down. Jax's trace program is eating through your ghost nodes.
[HACKING CHALLENGE]
Lyra: "I've queued up three scripts. One will download the dossier and wipe the logs simultaneously. The others... are too slow or too loud. Choose now!"`,
			Choices: []Choice{
				{Text: "run quick_download.bat /target Project_Chimera_Dossier", NextLocationKey: "game_over_jax_trace"},
				{Text: "run script_package.sh --action grab_and_wipe --target chimera_dossier.zip", NextLocationKey: "vengeance_chapter_2_success"},
				{Text: "run full_system_audit.exe", NextLocationKey: "game_over_jax_trace"},
			},
		},
		"vengeance_chapter_2_success": {
			Title: "Leverage Acquired",
			Description: `The script executes flawlessly. The dossier streams to your secure storage just as the final log entry of your session vanishes. Jax's trace program hits a dead end. You jack out, your heart pounding. You have it. The data proves Thorne's direct involvement in 'Project Chimera', and it names his superior: a Director named Evelyn Reed. You not only have blackmail on Thorne, you have your next target.

CHAPTER 2 COMPLETE. TO BE CONTINUED...`,
		},

		// --- PATH OF LIBERATION ---
		"liberation_intro": {
			Title:       "Path of Liberation: First Clue",
			Description: `"I knew there was a conscience still in you, Kael," Lyra says warmly. "Our first step is to find 'The Undercurrent,' a hidden network the Echoes use to communicate. The entry point is hidden in an old, abandoned subway station in Sector 7. Be careful. OmniCorp patrols are heavy there."`,
			Choices: []Choice{
				{Text: "Navigate the shadows to Sector 7.", NextLocationKey: "liberation_subway_approach"},
			},
		},
		"liberation_subway_approach": {
			Title:       "Sector 7 - Abandoned Subway",
			Description: `You arrive at the perimeter of the abandoned Sector 7 station. Two OmniCorp Enforcers stand guard. A flickering neon sign across the street casts long shadows. Lyra's voice materializes in your ear. "The Enforcers' comms are routed through that public terminal across the street. A little 'network interference' could create the distraction we need."`,
			Choices: []Choice{
				{Text: "Hack the public terminal to create a diversion.", NextLocationKey: "liberation_hack_distraction"},
				{Text: "Try to find another way in.", NextLocationKey: "liberation_alternate_route"},
			},
		},
		"liberation_hack_distraction": {
			Title: "Hacking Minigame: Network Interference",
			Description: `You jack into the terminal. Inside the local network, you see the Enforcers' comms channel and a link to the district's 'Civic Alert' system.
[HACKING CHALLENGE]
Lyra: "We need to flood their comms with noise. Rerouting the Civic Alert audio feed should do it. I've highlighted the primary node. Execute the command."`,
			Choices: []Choice{
				{Text: "execute_command('reroute_audio -target comms.enforcer_s7 -source system.civic_alert')", NextLocationKey: "liberation_undercurrent_entrance"},
				{Text: "execute_command('shutdown -target comms.enforcer_s7')", NextLocationKey: "game_over_heist_caught"},
			},
		},
		"liberation_alternate_route": {
			Title:       "Alternate Route",
			Description: `You scout the perimeter and find a collapsed section of wall leading into a maintenance tunnel. The air is thick with mold, but it appears unguarded. It's a risk, but it avoids a direct confrontation.`,
			Choices: []Choice{
				{Text: "Enter the maintenance tunnel.", NextLocationKey: "liberation_undercurrent_entrance"},
			},
		},
		"liberation_undercurrent_entrance": {
			Title:       "The Undercurrent",
			Description: `You slip past the distracted or oblivious guards and descend into the station. Following Lyra's directions, you find a rusted maintenance panel. Behind it, a fiber-optic cable, pulsing with a soft, cyan light, is spliced into the station's old infrastructure. This is the entrance to 'The Undercurrent'.`,
			Choices: []Choice{
				{Text: "Connect to The Undercurrent.", NextLocationKey: "liberation_chapter_end"},
			},
		},
		"liberation_chapter_end": {
			Title:       "Liberation: A New Hope",
			Description: `You jack in. The chaotic noise of the public Nexus fades, replaced by a quiet, orderly stream of encrypted data. You are in The Undercurrent. This hidden digital sanctuary is the nerve center of the resistance. Your fight for the soul of Veridia has just begun.`,
			Choices: []Choice{
				{Text: "Announce your presence in The Undercurrent.", NextLocationKey: "liberation_chapter_2_start"},
			},
		},
		"liberation_chapter_2_start": {
			Title: "Chapter 2: The Undercurrent",
			Description: `Your presence is immediately flagged. A text string scrolls across your vision, stark and unadorned.
[USER: Cipher] "You're new. And Lyra vouches for you. But you're ex-OmniCorp, which makes you a liability. Words are cheap. We need to see if your skills are real."`,
			Choices: []Choice{
				{Text: `"I'm ready. What's the test?"`, NextLocationKey: "liberation_mission_briefing"},
			},
		},
		"liberation_mission_briefing": {
			Title: "The Test Mission",
			Description: `[USER: Cipher] "An OmniCorp data analyst, name of Anya, wants to defect. She has information on 'Project Chimera'. She's being held for 'psychological re-calibration' at the Nightingale Clinic in the mid-levels. Security is tight, but we have a window. Get her out."
Lyra adds, "We have two viable infiltration vectors, Kael. We can try a remote hack to disable their security systems, or you can go in person and bluff your way past the front desk as a maintenance tech."`,
			Choices: []Choice{
				{Text: "[Digital Infiltration] Hack the clinic's security systems remotely.", NextLocationKey: "liberation_clinic_digital_approach"},
				{Text: "[Social Engineering] Go to the clinic and pose as a technician.", NextLocationKey: "liberation_clinic_social_approach"},
			},
		},
		"liberation_clinic_digital_approach": {
			Title: "Hacking Minigame: Digital Infiltration",
			Description: `You find a secure line into the clinic's network. It's an older system, but robust. Lyra highlights two primary control systems.
[HACKING CHALLENGE]
Lyra: "You can either create a cascading reboot of the camera and lock systems, causing a temporary but chaotic shutdown, or you can try to insert a 'ghost' credential into the active security roster, making the system think you belong there. The first is noisy, the second is subtle but complex."`,
			Choices: []Choice{
				{Text: "initiate --system_reboot --all", NextLocationKey: "liberation_clinic_success"},
				{Text: "inject --credential ghost.tech --roster active_security", NextLocationKey: "game_over_clinic_alarm"},
			},
		},
		"liberation_clinic_social_approach": {
			Title:       "Social Engineering: The Front Desk",
			Description: `You arrive at the sterile, white entrance of the Nightingale Clinic, wearing a scavenged technician's jumpsuit. The receptionist, a bored-looking man with a chrome optic, looks up. "ID and work order," he says flatly. You don't have a valid work order. This is all about the bluff.`,
			Choices: []Choice{
				{Text: `"There's a priority-one network diagnostic. Didn't they tell you?"`, NextLocationKey: "liberation_clinic_success"},
				{Text: `"I'm here to fix the... uh... hydro-spanners."`, NextLocationKey: "game_over_bluff_failed"},
			},
		},
		"liberation_clinic_success": {
			Title:       "The Extraction",
			Description: `Your plan works. The confusion from the system reboot or the confidence of your bluff gets you inside. You locate Anya in a holding room, dazed but unharmed. You escort her out through a service exit and into the anonymous crowds of the mid-level. You've passed the test.`,
			Choices: []Choice{
				{Text: "Return to The Undercurrent with Anya.", NextLocationKey: "liberation_chapter_2_success"},
			},
		},
		"liberation_chapter_2_success": {
			Title: "Trust Earned",
			Description: `Back in the digital sanctuary of The Undercurrent, Cipher's tone is different.
[USER: Cipher] "You did good. Very good. Anya's data confirms our worst fears about Chimera. Welcome to the resistance, Kael. The real fight starts now."
You have proven your worth. You are no longer just an outcast; you are one of the Echoes.

CHAPTER 2 COMPLETE. TO BE CONTINUED...`,
		},

		// --- PATH OF SELF-PRESERVATION ---
		"preservation_intro": {
			Title:       "Path of Self-Preservation: The Black Market",
			Description: `"You're making a mistake," Lyra warns, her voice laced with disappointment. "This information is too powerful." You ignore her, wiping her from your system. You put out feelers on the shadow nets. The highest bidder is a fixer named 'Silas'. The meet is at a noodle stand in the rain-drenched 'Gutter's Market'. You can't shake the feeling you're being watched.`,
			Choices: []Choice{
				{Text: "Go meet Silas.", NextLocationKey: "preservation_meet_silas"},
			},
		},
		"preservation_meet_silas": {
			Title:       "Gutter's Market",
			Description: `Silas is a rail-thin man with a nervous twitch and cybernetic eyes that scan the crowds constantly. "You got the data?" he asks, not making eye contact. "The client is... particular. They want proof before the transfer. Show me a snippet. The transaction will happen on my mark."`,
			Choices: []Choice{
				{Text: "Show him a meaningless, encrypted snippet.", NextLocationKey: "preservation_deal_success"},
				{Text: "Refuse. Demand the credits first.", NextLocationKey: "game_over_betrayed"},
			},
		},
		"preservation_deal_success": {
			Title:       "Self-Preservation: The Payout",
			Description: `The deal with Silas goes surprisingly smooth. You walk away with a fortune in untraceable credits. As you melt back into the crowd, your personal comm buzzes with a high-priority alert: your old OmniCorp security credentials have been flagged for termination by 'an interested party'. They're not just firing you; they're erasing you. Silas's client is tying up loose ends. You need to disappear, now.`,
			Choices: []Choice{
				{Text: "Find a way off-world. Immediately.", NextLocationKey: "preservation_chapter_2_start"},
			},
		},
		"preservation_chapter_2_start": {
			Title:       "Chapter 2: The Ghost Broker",
			Description: `You can't get off-world with a flagged identity. You need a ghost ID, and there's only one broker in the Neon Slums with the skill to craft one that will pass an orbital checkpoint: a woman known as 'Fade'. You find her operating out of a cramped cybernetics den, surrounded by flickering monitors displaying stolen identities.`,
			Choices: []Choice{
				{Text: `"I need a new face. The best you've got."`, NextLocationKey: "preservation_identity_choice"},
			},
		},
		"preservation_identity_choice": {
			Title:       "A Costly Face",
			Description: `Fade looks you over with unsettlingly calm eyes. "I can give you two kinds of new," she says, her voice a low hum. "I can give you a top-tier 'Platinum' identity. Custom biometrics, deep-level corporate back-stopping. It'll pass any scan they throw at you. It will cost you 80% of your payout. Or... I can give you a recycled 'Bronze' ID. It'll get you through the gate... probably. But it might have some old ghosts attached. Your call."`,
			Choices: []Choice{
				{Text: "Pay for the Platinum ID. Security is worth the price.", NextLocationKey: "preservation_checkpoint_premium"},
				{Text: "Gamble on the Bronze ID. I need the extra cash.", NextLocationKey: "preservation_checkpoint_cheap"},
			},
		},
		"preservation_checkpoint_premium": {
			Title:       "Veridia Orbital Port: First Class",
			Description: `You hand over the credits. Fade gets to work. An hour later, you are 'Alex Mercer', a mid-level logistics manager with a clean record. At the spaceport, the OmniCorp checkpoint guard waves you through without a second glance. The scan is green. You walk onto the transport, leaving the rain-slicked streets of Veridia behind for good.`,
			Choices: []Choice{
				{Text: "Watch Veridia disappear beneath you.", NextLocationKey: "preservation_chapter_2_success"},
			},
		},
		"preservation_checkpoint_cheap": {
			Title:       "Veridia Orbital Port: A Tense Wait",
			Description: `You take the cheaper ID. You're now 'Riktor', a cargo hauler with a slightly suspicious travel history. At the spaceport checkpoint, the gate flashes yellow as you pass through. The OmniCorp guard gestures you to a secondary screening. "Just a random check, sir," he says, his eyes cold. He scans your ID again. An alert flashes on his screen. "This ID was flagged for a debt default three hours ago. You have anything to say about that, Riktor?"`,
			Choices: []Choice{
				{Text: `"That's a mistake. Run it again."`, NextLocationKey: "game_over_checkpoint_fail"},
				{Text: "[BRIBE] Slip him a cred-chip. \"Maybe your scanner needs a 'calibration'.\"", NextLocationKey: "preservation_chapter_2_success"},
			},
		},
		"preservation_chapter_2_success": {
			Title: "The Escape",
			Description: `Whether through flawless preparation or a well-placed bribe, you make it onto the transport. As the ship breaks atmosphere, you look down at the sprawling neon grid of Veridia. You are free. You are rich. But as you watch the city shrink into a pinpoint of light, you can't shake the feeling that the client who bought your data is still out there, and they now know what you look like.

CHAPTER 2 COMPLETE. TO BE CONTINUED...`,
		},

		// --- GAME OVER STATES ---
		"game_over_static": {
			Title:       "GAME OVER: The Static Life",
			Description: `You leave the strange fragment behind. The paydata you fence is enough for a few weeks' rent and nutrient paste. The days blur into one another, a monotonous cycle of small-time hacks and paranoid glances over your shoulder. You never learn the truth. You remain a ghost, a nobody, another piece of static lost in the noise of Veridia.`,
		},
		"game_over_heist_caught": {
			Title:       "GAME OVER: T-Minus Ten Seconds",
			Description: `Your command is too aggressive. A silent alarm triggers. An Intrusion Countermeasure daemon flashes to life, tracing your connection in milliseconds. Before you can react, your screen goes black. Ten seconds later, your door splinters inward as OmniCorp Enforcers storm your location. There is nowhere to run.`,
		},
		"vengeance_social_fail": {
			Title:       "GAME OVER: Amateur Hour",
			Description: `You approach Thorne, trying to spin a story. He isn't a supervisor for nothing. He sees through your amateur attempt at social engineering in seconds. "Security," he says calmly into his wrist-comm, never taking his eyes off you. "We have a Static problem." You don't make it two steps before you are apprehended.`,
		},
		"game_over_betrayed": {
			Title:       "GAME OVER: A Knife in the Dark",
			Description: `You scoff. "Credits first." Silas almost looks sad. "Bad choice," he whispers. You feel a sharp, cold pain in your side. One of his associates stuck you with a neuro-toxin injector in the crowd. Your vision blurs as you slump over the noodle counter. Your last sight is Silas lifting your deck from your jacket. You gambled and lost.`,
		},
		"game_over_clinic_alarm": {
			Title:       "GAME OVER: The Silent Alarm",
			Description: `Your command is too elegant for this old system. It's flagged as an unknown intrusion. You don't see the alarm, but a security team is dispatched to the network access point you're using. They find you still jacked in, an easy target. Your fight ends before it truly began.`,
		},
		"game_over_bluff_failed": {
			Title:       "GAME OVER: A Poor Lie",
			Description: `The receptionist's chrome optic whirs as it scans you. "Hydro-spanners? We don't have those." He presses a button under his desk. "Security, we have a Static trying to gain entry." Two heavily armed OmniCorp Enforcers emerge from a side door. There is no escape.`,
		},
		"game_over_checkpoint_fail": {
			Title:       "GAME OVER: Loose Ends",
			Description: `You try to bluff. The guard isn't buying it. He calmly speaks into his comm. "We have a runner." The port's blast doors slam shut. Floodlights pin you in place. The client who bought your data was thorough. You were the last loose end, and now you've been tied up permanently.`,
		},
		"game_over_jax_trace": {
			Title: "GAME OVER: The Ghost is Caught",
			Description: `You chose the wrong script. It's too slow and noisy. Jax's trace locks onto your physical location instantly. "Trace complete," a system message flashes. "Enforcers dispatched."
Lyra's voice is a pained whisper. "I'm sorry, Kael..." The connection dies. You don't have time to run.`,
		},
	}

	return &Game{
		Player:    player,
		Locations: locations,
		IsRunning: true,
		Logger:    logger,
	}
}

// --- Game Logic ---

// Run starts and manages the main game loop.
// It uses a goroutine to handle user input concurrently, allowing the main loop
// to remain responsive.
func (g *Game) Run() {
	g.Logger.Info("Veridian Echo game loop started.")

	inputChan := make(chan string)
	quitChan := make(chan struct{})

	// Per the Concurrency Mandate, I/O is handled in a separate goroutine.
	go g.readInput(inputChan, quitChan)

	// Initial render of the starting location.
	g.renderCurrentLocation()

gameLoop:
	for {
		select {
		case <-quitChan:
			// Input stream closed (e.g., Ctrl+D), terminate game.
			g.IsRunning = false
			break gameLoop
		case input := <-inputChan:
			// Process received input.
			g.handleInput(input)
			if !g.IsRunning {
				// The game has ended based on the last choice.
				break gameLoop
			}
			g.renderCurrentLocation()
		}
	}

	g.Logger.Info("Veridian Echo game loop terminated.")
	fmt.Println("\n[CONNECTION_TERMINATED]")
}

// readInput runs in a dedicated goroutine, listening for user input
// without blocking the main application loop.
func (g *Game) readInput(inputChan chan<- string, quitChan chan<- struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inputChan <- strings.TrimSpace(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		g.Logger.Error("Error reading from stdin", "error", err)
	}

	// Signal the main loop to quit if the scanner stops.
	close(quitChan)
}

// renderCurrentLocation displays the current scene's information to the player.
func (g *Game) renderCurrentLocation() {
	location, exists := g.Locations[g.Player.CurrentLocationKey]
	if !exists {
		g.Logger.Error("Player location key does not exist", "key", g.Player.CurrentLocationKey)
		g.IsRunning = false
		return
	}

	clearScreen()
	typewriteEffect(fmt.Sprintf("--- %s ---\n", location.Title), 40*time.Millisecond)
	fmt.Println()
	typewriteEffect(location.Description, 20*time.Millisecond)
	fmt.Println()

	// If a location has no choices, it's an end state.
	if len(location.Choices) == 0 {
		g.IsRunning = false
		return
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("--- CHOICES ---")
	for i, choice := range location.Choices {
		fmt.Printf("%d: %s\n", i+1, choice.Text)
	}
	fmt.Println("---------------")
	fmt.Print("> ")
}

// handleInput parses and validates the user's choice, updating the game state.
func (g *Game) handleInput(input string) {
	location := g.Locations[g.Player.CurrentLocationKey]
	choiceCount := len(location.Choices)

	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("\nInvalid command. Please enter a number.")
		time.Sleep(1 * time.Second)
		return
	}

	if choice < 1 || choice > choiceCount {
		fmt.Printf("\nInvalid choice. Please enter a number between 1 and %d.\n", choiceCount)
		time.Sleep(1 * time.Second)
		return
	}

	chosen := location.Choices[choice-1]

	// Check for and execute a state-changing action if one exists.
	if chosen.Action != nil {
		chosen.Action(g.Player)
	}

	g.Player.CurrentLocationKey = chosen.NextLocationKey

	g.Logger.Info("Player made a choice",
		"choice_text", chosen.Text,
		"next_location", chosen.NextLocationKey,
		"player_path", g.Player.Path,
	)
}

// --- Utility Functions ---

// typewriteEffect prints text to the console with a delay between characters
// for atmospheric effect.
func typewriteEffect(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(delay)
	}
}

// clearScreen clears the console.
// This is a simple, OS-agnostic implementation that prints newlines.
// For a production executable, OS-specific clear commands would be used,
// but this maintains single-file simplicity.
func clearScreen() {
	fmt.Print(strings.Repeat("\n", 50))
}

// --- Main Execution ---

func main() {
	// Per Pillar IV, initialize structured JSON logging.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Per Pillar II, initialize the game state from our content definition.
	game := initializeGame(logger)

	// Run the game's core loop.
	game.Run()
}
