# Data Model: Specify Progress Feedback

## Entity: ExecutionPhase

- **Description**: Named stage of `/speckit-specify` command execution.
- **Key Fields**:
  - `phase_id`
  - `phase_name`
  - `phase_type` (`hook`, `generation`, `validation`, `completion`, `failure`)
  - `started_at`
  - `ended_at`
- **Validation Rules**:
  - Every major phase must emit at least one visible progress event.
  - Terminal phase states must be mutually exclusive (`completed` or `failed`).

## Entity: ProgressEvent

- **Description**: User-visible status signal describing runtime command state.
- **Key Fields**:
  - `event_type` (`started`, `heartbeat`, `completed`, `skipped`, `failed`)
  - `phase_id`
  - `message`
  - `timestamp`
  - `reason` (required for `skipped` and `failed`)
- **Validation Rules**:
  - `heartbeat` events appear at least every 10 seconds for long-running phases.
  - `skipped` events must include step identity and skip reason.

## Entity: RenderingMode

- **Description**: Output rendering path selected for progress feedback.
- **Key Fields**:
  - `mode` (`rich`, `plain_text`)
  - `fallback_trigger`
- **Validation Rules**:
  - On rich rendering failure, mode must transition to `plain_text` without aborting execution.

## Relationships

- `ExecutionPhase` produces one or more `ProgressEvent` records.
- `RenderingMode` controls presentation format of each `ProgressEvent`.
