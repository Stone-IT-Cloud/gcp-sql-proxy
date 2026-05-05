# Quickstart: `/speckit-specify` Progress Feedback

## 1. Validate phase updates

Run `/speckit-specify` in a normal flow.

Expected:
- Visible phase progress updates appear for major command stages.
- Final output includes completion confirmation and generated artifact summary.

## 2. Validate hook progress transitions

Run with pre-execution hooks enabled.

Expected:
- Hook start and hook completion events are shown explicitly.

## 3. Validate long-running liveness cadence

Trigger a longer-running planning/generation step.

Expected:
- Liveness updates are emitted at least every 10 seconds until completion/failure.

## 4. Validate skipped optional step behavior

Execute scenario where an optional step is skipped.

Expected:
- A `skipped` event appears including the skipped step name and reason.

## 5. Validate failure context

Force an execution error in a known phase.

Expected:
- Output identifies the failed phase.
- Output includes at least one actionable next step.

## 6. Validate rendering fallback

Simulate rich renderer failure.

Expected:
- Output falls back to plain text.
- Command continues running to completion/failure.
