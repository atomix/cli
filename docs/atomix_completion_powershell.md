## atomix completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	atomix completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
atomix completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [atomix completion](atomix_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 8-May-2022
