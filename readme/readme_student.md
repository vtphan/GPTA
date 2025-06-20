### Student's installation

#### 1. Install Sublime Text
First, download and install Sublime Text if you haven’t already. (https://www.sublimetext.com/download)


#### 2. Install the GEMStudent Plugin
Step 1: Open Sublime Text
    •	Launch Sublime Text on your computer.

Step 2: Open the Console
    •	Go to the View menu and select Show Console.

Step 3: Installation Command
    Copy the following command, paste it into the console in sublime text and hit enter:

**import os; package_path = os.path.join(sublime.packages_path(), "GEMStudent"); os.mkdir(package_path) if not os.path.isdir(package_path) else print("dir exists"); module_file = os.path.join(package_path, "GEMStudent.py") ; menu_file = os.path.join(package_path, "Main.sublime-menu"); version_file = os.path.join(package_path, "version.go"); import urllib.request; urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GPTA/2.1/src/GEMStudent/GEMStudent.py", module_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GPTA/2.1/src/GEMStudent/Main.sublime-menu", menu_file); urllib.request.urlretrieve("https://raw.githubusercontent.com/vtphan/GPTA/2.1/src/version.go", version_file)**
    
Step 4: Restart Sublime Text



#### 3. Configure the Plugin: 
configuration settings will be provided by your instructor.

![alt text](</assets/documentation/image2.png>)


 follow these steps:

    1. Set Server Address

    2. Set course id

    3. Set username

    4. login

Once successfully logged in, GEM plugin menu should look like this.
![alt text](/assets/documentation/image1.png)
