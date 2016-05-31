# (Free)SwitchMon(itor) Tool

A tool to get info out of Freeswitch and send it to Riemann

## Dependencies

  $ make deps

## Build

  $ make

  That will create a `switch_mon` binary under the `release` folder.

## Usage

  Copy `switch_mon` binary to an instance where Freeswitch is running. You'll need to be running the EventSocket module.

  Create a `config.json` within the same bin folder:

  An example `config.json` would be:

```
{
  "hooks": [
    {
      "event": "CHANNEL_ANSWER",
      "attributes": ["variable_sip_to_uri"],
      "service": "answer"
    },
    {
      "event": "CHANNEL_HANGUP_COMPLETE",
      "metrics": ["variable_rtp_audio_in_quality_percentage", "variable_rtp_audio_in_mos", "variable_rtp_audio_in_quality_percentage", "variable_rtp_audio_in_mos"],
      "attributes": ["variable_sip_to_uri", "variable_sip_user_agent"],
      "service": "hangup"
    }
  ]
}
```

  Finally, run it with:

  `switch_mon -riemann-host your-riemann-host.com`

  See `switch_mon -h` for extra args.
