# World of Warcraft - Quest Reader

This WoW Addon will read the quest text with meant of AI text-to-speech API(s).

# Preconditions (what do you need)

* [AWS Polly](https://aws.amazon.com/polly/) credentials

# Installation

1. Unzip the [WoW_Addon.zip](https://github.com/rainu/wow-quest-reader/releases/) to 
    `&lt;World of Warcraft&gt;\_retail_\Interface\AddOns`
2. Download the [companion-windows-amd64.exe](https://github.com/rainu/wow-quest-reader/releases/)-Application and move them into a directory you want 
3. Create a file named `config.yml` **in the same directory** as the Companion-Application
   1. Insert at least the AWS-Credentials
   2. For more options see the **configuration** section
4. Start the Companion-Application
5. Start World of Warcraft
   1. Go to Keybinding settings in WoW
   2. Set the key for **Rainu Quest Reader** &gt; **Collect last quest information** to [ctrl]+&lt;F12&gt; 
      (or what ever you have configured - see **key.addon** in configuration section!)
6. Now if you open a quests you can press the read button (&lt;PAGE DOWN&gt; - or what ever you configured - see **key.read** in configuration section!)

# Configuration

| Configuration | Default | Mandatory | Description |
|---|---|---|---|
| debug | false | false | Is the application running in Debug-Mode? |
| logLevel | 4 | false | The used logging level. ( Panic(0); Fatal(1); Error(2); Warn(3); Info(4); Debug(5);Trace(6) ) |
| sound.directory | &lt;Companion-Application&gt;/sounds | false | The directory where the generated sound files will be stored. |
| sound.aws.region |  | true | The AWS region which should be used. |
| sound.aws.key |  | true | The AWS key which should be used. |
| sound.aws.secret |  | true | The AWS secret which should be used. |
| key.read | &lt;PAGE DOWN&gt; | false | The keybinding when the application should start reading. |
| key.addon | [ctrl]+&lt;F12&gt; | false | The Keybinding for the addon. |

## Example config.yml

```yaml
key:
  read: "[ctrl]+<F12>"
sound:
  dir: /tmp/wow
  aws:
    region: "<region>"
    key: "<key>"
    secret: "<secret>"
```

# How it works

The Companion-Application observe the clipboard. The addon will collect the quest data while you playing WoW. 
If the shortcut (**key.addon**) is pressed, the addon opens a hidden window within all quest information as text selected.
Then the Companion-Application will copy the selected content into the clipboard by "pressing" &lt;CTRL&gt;+&lt;C&gt;. 
The quest data is copied into the clipboard and the application will detect that content and will generate and
play the speech. Because the generation will potential cost some money, the sound file fill be persisted inside the 
configured sound folder. So that play the same quest multiple times will not cost for every playing.

# Build

```shell
GOOS=windows GOARCH=amd64 go build -a -o companion.exe ./app/companion
```