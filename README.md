# next-task

This project is meant to experiment with pulling together tasks from multiple
sources and selecting which one work on next automatically. The project is just
for fun and for self use so it may be a little hacky in some spots.

This will require a prioritized backlog of tasks that each range from 15 minutes
to 3 hours. The tool will load a block of the top 20 tasks and pick one at random
to work on until that block is complete.

The size of the tasks and randomness is part of the experiment. I
want to see if small tasks in no particular order can help ensure more projects
are worked on in parallel so that a larger project doesn't end up stalling
other ones out.

The reason a chunk of tasks are chosen from a prioritized list is also to make sure
that things aren't too random and the important work is still being worked on.

**Why 20?**

15m * 20 = 5h which is an average working day sans meetings, ad-hoc requests,
and other time bound items. And if I ever wanted do implement this on paper, I
could do it with a d20 die.

## Status

Just getting started on it. Not battle tested at all yet.

## Supported Sources

### Google Tasks and Gmail

This tool requires a Google project set up to generate the
`~/.config/next-task/google/credentials.json` file, and the project should have
read access to Google Tasks and Gmail. The quickstarts for both of these
services should provide enough access.

## Planned Sources

### Jira

A token and JQL statement to select which tasks are in scope will be needed
