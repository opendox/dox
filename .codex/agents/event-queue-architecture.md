<!--
  dox
  Copyright (C) 2026  OpenDox

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program. If not, see <http://www.gnu.org/licenses/>.

  @File    : .codex/agents/event-queue-architecture.md
  @Author  : Frost Leo <frostleo.dev@gmail.com>
  @Created : 2026-04-18
  @Modified: 2026-04-18
-->

# 6. Event And Queue Architecture


Dox systems should collaborate through explicit tasks and events.

Web, Scheduling, Collection, and Computation must not become tightly coupled through hidden in-process calls when work crosses system boundaries. Work that is cross-system, long-running, retryable, traceable, or potentially distributed should be coordinated through task records and queue messages.

Events and queues are not added for complexity. They exist so Dox can keep clear boundaries, traceability, recoverability, and scalability across multiple processes, machines, asynchronous tasks, large-scale collection, and computation workflows.

### Event-driven Principles

Important state changes in Dox should become events where practical.

Typical events include:

- Plugin configuration changed.
- Plugin enabled or disabled.
- Platform credential changed.
- Collection task created.
- Collection task started.
- Collection task completed.
- Collection task failed.
- Collection result persisted.
- Computation task created.
- Computation task started.
- Computation task completed.
- Computation task failed.
- Computation result persisted.
- System task alert triggered.
- Business alert input generated.
- Notification sent.

Events should include traceability fields where practical:

- `event_id`
- `event_type`
- `task_id`
- `plugin_id`
- `platform_id`
- `source_id`
- reserved `tenant` or `workspace` fields
- `correlation_id`
- `causation_id`
- `occurred_at`

Even if strict multi-tenancy is not implemented in the current project scope, key events and tables may reserve names and boundaries needed for future extension. Do not overbuild the current implementation around future multi-tenancy.

### Tasks And Events

Tasks represent work the system needs to perform.

Events represent facts that have already happened.

For example:

- "Create a keyword collection task" is a task.
- "The keyword collection task completed" is an event.
- "Create a search result aggregation task" is a task.
- "Search result aggregation completed" is an event.

Do not mix task semantics and event semantics. Tasks may fail, retry, time out, or be cancelled. Events should represent facts that already happened and must not be treated as arbitrary commands.

### Queue Technology

Dox should prefer Kafka as the core queue and event-stream platform.

Kafka is preferred because Dox is expected to have:

- Large-scale collection events.
- Data persistence events.
- Task status events.
- Computation trigger events.
- Replayable data flows.
- Future report and computation inputs.
- Event collaboration across many workers and systems.

Kafka is a better long-term foundation for event logs and data pipelines than a simple work queue.

Dox must not introduce Kafka and RabbitMQ together without clear boundaries. RabbitMQ or another queue may be added later only with a documented responsibility, such as delayed tasks, complex routing, work queues, or specific RPC-like scenarios.

If a second queue technology is introduced later, its responsibility must be clearly separated from Kafka. The same task category must not exist under two conflicting queue semantics.

### Scheduling And Queues

The Scheduling System is the task decision-maker.

The Scheduling System creates tasks, decides when tasks enter queues, decides task priority, manages task state, and advances the next step based on results.

The Collection System and Computation System are task executors.

They may consume queue messages, execute tasks, write results, report status, and emit events. They must not own global task chains or scheduling policy.

### Queue And Task Operations

Dox must provide observability and intervention capabilities for queues and tasks.

The Web System should provide management entries that allow authorized users or operators to inspect and manage scheduling queues, collection task queues, and computation task queues.

Queue and task operations may include:

- Viewing queue status.
- Viewing task lists.
- Viewing task details.
- Viewing task priority.
- Viewing task dependencies.
- Viewing task execution history.
- Viewing task failure reasons.
- Viewing task retry counts.
- Viewing related plugin, platform, configuration, and data source information.
- Pausing tasks.
- Resuming tasks.
- Cancelling tasks.
- Re-running tasks.
- Manually triggering tasks.
- Adjusting task priority.
- Reordering tasks.
- Changing scheduled execution time.
- Moving tasks into or out of failed-task handling.
- Batch retrying failed tasks.
- Adding manual notes or operator marks to abnormal tasks.

These operations must be protected by authorization, audit logging, and safety boundaries. Ordinary users must not be able to modify system-level queues or tasks owned by other users.

Queue management must not bypass or corrupt the task state machine. Any manual intervention must be represented through explicit state transitions, audit logs, and scheduling rules.

If Kafka is used underneath, Kafka topics must not be treated like database tables that can be freely updated, reordered, or edited. Dox should maintain manageable task state, task plans, and task indexes in the Scheduling System. The Web System operates on the Dox task model and scheduling model, not directly on the underlying message log.

Queue messages are execution channels. The task model is the operations management object.

### Idempotency And Recovery

All task handlers and event handlers that may be consumed more than once must be designed for idempotency.

The system must assume:

- Queue messages may be duplicated.
- Workers may crash during execution.
- Tasks may fail halfway.
- Status callbacks may be lost.
- External APIs may rate-limit, time out, or return unstable results.
- The same data may be collected more than once.

Task and event handling should define:

- Idempotency keys.
- Deduplication keys.
- State machines.
- Retry policy.
- Timeout policy.
- Dead-letter or failure records.
- Traceable logs.
- Task execution records.
- Relationships between raw data and processed results.

### Event Naming And Documentation

Event names must be clear, stable, and readable.

Event names should describe facts that already happened, such as:

- `plugin.enabled`
- `collection.task.created`
- `collection.task.completed`
- `collection.task.failed`
- `collection.result.persisted`
- `computation.task.created`
- `computation.task.completed`
- `computation.result.persisted`
- `alert.system.triggered`
- `alert.business.input.generated`
- `notification.sent`

Event structure and semantics should be documented in code comments or architecture documents. Events used across systems must not change fields casually. Use event versions when necessary.

