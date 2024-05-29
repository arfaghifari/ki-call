from os import walk


class ClientFile :
    def __init__(self, module, package, path, filename, proto_path):
        self.module = module
        self.package = package
        self.path = path
        self.lines = ["// This Code is Generated"]
        self.filename = path+ "/" +filename
        self.lib = []
        self.objs = []
        self.proto_path = proto_path

    def compile_lines(self):
        self.add_package_line()
        self.add_empty_line(1)
        self.add_lib_line()
        self.add_empty_line(1)
        self.add_vairabel_line()
        self.add_empty_line(1)
        self.add_struct_line()
        self.add_empty_line(1)
        self.add_func_line()

    def generate(self):
        f = open(self.filename, "w")
        for line in self.lines:
            f.write(line + "\n")
        f.close()

    def add_empty_line(self, n):
        for _ in range(n):
            self.lines.append("")
    
    def add_package_line(self):
        self.lines.append("package" + " "+ self.package)
    
    def add_lib(self, lib):
        self.lib.append(lib)
    
    def add_lib_line(self):
        if len(self.lib) <= 0 :
            return
        self.lines.append("import(")
        for lib in self.lib :
            self.lines.append("\t" +"\"" +lib + "\"" )
        self.lines.append("\t\"github.com/cloudwego/kitex/client\"")
        self.lines.append(")")

    def add_obj(self, obj):
        self.objs.append(obj)
        self.add_lib(self.module + "/kitex_gen/" + obj + "/"+ obj)

    def generate_obj_from_proto_dir(self):
        for (_,_,filename) in walk(self.proto_path):
            for item in filename:
                self.add_obj(item.split(".")[0])

    
    def add_struct_line(self):
        self.lines.append("type Client struct {")
        for obj in self.objs:
            objC = obj.capitalize()
            self.lines.append("\t" + objC + "  " + obj + ".Client")
        self.lines.append("}")

    def add_func_line(self):
        self.lines.append("func (c *Client) RegisterAllClient(host string) {")
        for obj in self.objs:
            objC = obj.capitalize()
            self.lines.append("\tClientKitex." + objC + ", _ = " + obj+".NewClient(\""+ obj + "\", client.WithHostPorts(host))")
        self.lines.append("}")

    def add_vairabel_line(self):
        self.lines.append("var ClientKitex Client")

    


CF = ClientFile("github.com/arfaghifari/ki-call","client", "src/client","client.go","proto")
CF.generate_obj_from_proto_dir()
CF.compile_lines()
CF.generate()