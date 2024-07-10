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

### Arcade

- presets/drs.json: DANCERUSH STARDOM
- presets/resident.json: beatmania IIDX 30 RESIDENT ~ current
- presets/vividwave.json: SOUND VOLTEX VIVID WAVE ~ current
- presets/around.json: DANCE aROUND

### コナステ/Konasute:

- presets/eac/bonga.json: Bombergirl
- presets/eac/ddr.json: DanceDanceRevolution GRAND PRIX
- presets/eac/gitadora.json: GITADORA
- presets/infinitas_2020.json: beatmania IIDX INFINITAS 2020
- presets/infinitas_2020_cache.json: beatmania IIDX INFINITAS 2020 (cache)
- presets/infinitas_2024.json: beatmania IIDX INFINITAS 2024
- presets/infinitas_2024_cache.json: beatmania IIDX INFINITAS 2024 (cache)
- presets/eac/nost.json: NOSTALGIA
- presets/eac/popn.json: pop'n music Lively
- presets/eac/sdvx.json: SOUND VOLTEX EXCEED GEAR
