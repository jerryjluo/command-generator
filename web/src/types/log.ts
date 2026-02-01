export interface TmuxInfo {
  in_tmux: boolean;
  session?: string;
  window?: string;
  pane?: string;
}

export interface ModelInput {
  system_prompt: string;
  user_prompt: string;
}

export interface ModelOutput {
  raw_response: string;
  command: string;
  explanation: string;
}

export interface Iteration {
  feedback: string;
  model_input: ModelInput;
  model_output: ModelOutput;
  timestamp: string;
}

export interface ContextSources {
  claude_md_content: string;
  terminal_context: string;
}

export interface Metadata {
  timestamp: string;
  model: string;
  final_status: 'accepted' | 'rejected' | 'quit';
  final_feedback?: string;
  iteration_count: number;
  tmux_info: TmuxInfo;
}

export interface SessionLog {
  id: string;
  user_query: string;
  context_sources: ContextSources;
  iterations: Iteration[];
  metadata: Metadata;
}

export interface LogSummary {
  id: string;
  user_query: string;
  final_status: 'accepted' | 'rejected' | 'quit';
  model: string;
  timestamp: string;
  iteration_count: number;
  command_preview: string;
  tmux_session?: string;
}

export interface LogListResponse {
  logs: LogSummary[];
  total: number;
  limit: number;
  offset: number;
}

export interface FilterParams {
  status?: string;
  model?: string;
  search?: string;
  from?: string;
  to?: string;
  sort?: string;
  order?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}
