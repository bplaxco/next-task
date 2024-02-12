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

## Completing Tasks

This doesn't edit the source where it gets the task. You still have to mark the
task complete when you're done. This is intentional so that you can get your
next task recommendation and start working on it but if you're unable to
complete it in the time window you can come back to it again when it gets
randomly pulled in the next round.

## Status

Just getting started on it. Not battle tested at all yet.

## Supported Sources

### Google Tasks and Gmail

This tool requires a Google project set up to generate the
`~/.config/next-task/google/credentials.json` file, and the project should have
read access to Google Tasks and Gmail. The quickstarts for both of these
services should provide enough access.

If the credentials.json file isn't present then these sources won't be
enabled.

### Jira

This pulls issues from Jira to use. The following settings need to be set
to fetch tasks from Jira:

- `NEXT_TASK_JIRA_INSTANCE_URL`: The instance to connect to
- `NEXT_TASK_JIRA_ACCESS_TOKEN`: The personal access token to make requests with
- `NEXT_TASK_JIRA_TASK_JQL`: Return the set of tasks that should be in scope

## Planned Sources

### Slack (Maybe)

If all of the above end up being useful for me, I may round this out with the ability
to connect to slack for pulling things from the save for later feature.

## Resetting the cache

To remove the task cache run:

```sh
rm -rf ~/.cache/next-task/tasks
```

## Misc Ideas

- Weighted Randomness: While pure randomness may feel fair, assigning slightly
  higher probabilities to the top-priority tasks ensures the most important
  items still get more attention overall.

- Categorization: Consider creating separate queues for different categories of
  tasks (e.g., urgent, deep-focus work, quick wins).
