# Contract: `/speckit-specify` Progress Feedback

## Scope

Defines user-visible progress behavior for `/speckit-specify` execution.

## Phase Event Contract

1. Command MUST emit progress updates for each major phase.
2. Progress text MUST identify the current phase in user-readable language.
3. Pre-execution hooks MUST emit start and finish events.
4. Spec generation and checklist validation MUST emit begin and complete events.

## Liveness Contract

1. Long-running phases MUST emit liveness events at least every 10 seconds.
2. Liveness output MUST indicate command is still active.

## Skip/Failure Contract

1. Optional skipped phases MUST emit a `skipped` event with step name and reason.
2. Failures MUST emit failed phase + at least one actionable next step.
3. Failure output MUST avoid internal stack traces or sensitive values.

## Rendering Resilience Contract

1. Rich rendering failures MUST trigger plain-text fallback.
2. Plain-text fallback MUST continue execution without aborting command.

## Completion Contract

1. Completion output MUST summarize generated outputs.
2. Completion output MUST be clearly distinguishable from in-progress states.
