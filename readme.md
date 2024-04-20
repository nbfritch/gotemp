# gotemp
A cli tool for getting a hwmon temperature by name.
Made specifically for use with waybar.

## Usage
`gotemp -name <name> -warn <warn temp> -crit <crit temp> -input <specific input>`

In waybar's `config.json` file:
```
  {
    "modules-right": [
      <other-entries>
      "custom/temperature",
      <other-entries>
    ],
    "custom/temperature": {
      "format": "{}",
      "exec": "gotemp -name k10temp -warn 70 -crit 90",
      "interval": 15,
      "return-type": "json"
    }
  }
```

In waybar's `style.css` file:
```css
#custom-temperature {} /* Normal */
#custom-temperature.warning {} /* When above warn temp */
#custom-temperature.critical {} /* When above crit temp */
```

### Example
`gotemp -name k10temp -warn 75 -input 1`
Outputs
```json
{"text": "60", "class": "normal"}
```

## Options

### Required options
`-name` The name listed in `/sys/class/hwmon/hwmonN/name`

### Optional options
`-warn <T>` The temperature in celsius after which the warning class is enabled. Default 70
`-crit <T>` The temperature in celsius after which the critical class is eanbled. Default 90
`-verbose` Whether or not to enable verbose logging. (Debug only, interferes with waybar)
`-input` The specific input under the hwmon node to read from. Default behavior is to average the values of `temp*_input`

