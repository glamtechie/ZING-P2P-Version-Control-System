# P2P-Version-Control-System

Filesystem APIs:
We expose the following APIs to the layers above:
- func zing_init(id int): Initializes the zing repository
- func zing_pull(branch string): Pulls changes from the global repository to the given branch in the local repository 
- func zing_add(filename string): Adds changes for commiting 
- func zing_commit(): Commits changes added so far 
- func zing_make_patch_for_push(branch string, patchname string): 
  Creates patch file by the name provided for the commited changes from the local.  This patch file needs to be propagated to other nodes based on the 2PC protocol.
- func zing_abort_push(): TODO
- func zing_process_push(patchname string): Applies the patch to the main branch in the global log
