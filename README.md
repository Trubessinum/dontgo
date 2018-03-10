# dontgo
Learning Go with OpenGL

# Quick start
### 1. Get that nifty [gvm](https://github.com/moovweb/gvm):

```
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
```

### 2. Get the last Go available:

```
gvm install go1.4
gvm use go1.4
export GOROOT_BOOTSTRAP=$GOROOT
gvm install go1.9.4
gvm use go1.9.4
```

### 3. Get the package:
```
git clone https://github.com/trubessinum/dontgo
cd dontgo
```

### 4. Run the code as follows:

```
go run main.go shader.go
```

### 5. See the result:
#### You should get something similar to this:
![pyramid](https://i.imgur.com/Acopm0F.gif)
