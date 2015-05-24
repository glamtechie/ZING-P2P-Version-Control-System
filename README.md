# P2P-Version-Control-System

Filesystem APIs:
We expose the following APIs to the layers above:
- func zing_init(id int): Initializes the zing repository
- func zing_pull(branch string): Pulls changes from the global repository to the given branch in the local repository 
- func zing_add(filename string): Adds changes for commiting 
- func zing_commit(): Commits changes added so far 
- func zing_push(branch string, patchname string): 
  Pushes commited changes from the local to the global on a temp branch and creates a patch file from the name provided in the argument.
  This patch file needs to be propagates to other nodes based on the 2PC protocol.
