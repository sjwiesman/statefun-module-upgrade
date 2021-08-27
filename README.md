Stateful Functions Automated Module Converter
=======

Automated tool to upgrade Stateful Function `module.yaml` files from the modern format released in 3.1.

### How To Use

Either Download an official [release](https://github.com/sjwiesman/statefun-module-upgrade/releases) or build the tool locally.

Then simply run the tool over your existing `module.yaml` files. 

```bash
$ cat module.yaml | ./statefun-module-upgrade 
```

### Why Change the format?

StateFun applications consist of multiple configuration components, including remote function endpoints, along with ingress and egress definitions, defined in a YAML format. 
Stateful Functions 3.1 added a new structure that treats each StateFun component as a standalone YAML document.
Thus, a `module.yaml` file becomes simply a collection of components.

While this might seem like a minor cosmetic improvement, this change opens the door to more flexible configuration management options in future releases - such as managing each component as a custom K8s resource definition or even behind a REST API.
StateFun still supports the legacy `module.yaml` file for backward compatibility, but users are encouraged to upgrade. 

