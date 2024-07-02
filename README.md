# extdrm

A tool for decrypting extdrm encrypted filesystems

# Usage

```
extdrm utility

Usage:
  extdrm [command]

Available Commands:
  dump        Decrypt the contents of an extdrm encrypted filesystem
  help        Help about any command
  verify      Verify the integrity of a filesystem dump

Flags:
  -h, --help   help for extdrm

Use "extdrm [command] --help" for more information about a command.
```

To dump the contents of an encrypted filesystem, simply run `extdrm dump SOURCE DESTINATION PRESET`.

After a successful dump, you may perform a file check by running `extdrm verify DESTINATION`.

## Presets

This repository includes several preset files responsible for controlling the decryption process. These preset files can be found under the `presets` directory.

List of included preset files and their respective games:
- presets/drs.json: DANCERUSH STARDOM
- presets/resident.json: beatmania IIDX 30 RESIDENT ~ current
- presets/vividwave.json: SOUND VOLTEX VIVID WAVE ~ current
