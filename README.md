# kcauto

I wanted to watch Daria without dealing with ads or having to get up and
click the next episode on kisscartoon. So I made this playwright script to
partially solve that.

## Setup

Because of the ad atrocity on kisscartoon, I'd add a new firefox profile, 
then set the `FIREFOX_PROFILE` environment with the full path to the 
profile directory. Because playwright uses the nightly version of firefox,
you can use the option `setup-profile` for opening the browser and
adding the ublock extension.

```bash
firefox -CreateProfile <new-profile-name>
export FIREFOX_PROFILE=/full/path/to/firefox/profile/dir
kcauto --setup-profile

# ...
# Then you'll add ublock origin like any other firefox
# extension.
```

## Usage

You can either pass the url of the episode to start playing as an
argument, or use a 'last-episode' file that contains the url of the
last episode played. `kcauto` will automatically create one of these
unless the `-n/--no-write` option is specified.

```bash
usage: kcauto <url-to-episode>
	-l/--last-episode FILENAME	set last episode file to FILENAME
	-n/--no-write               do not create last-episode file
	-h/--help                   display this help message
	--setup-profile             open browser using firefox profile (for ublock setup)
```
