import { useEffect } from "react";
import { useAppStore } from "../stores/useAppStore";
import { EmptyState } from "../components/atoms";
import { ActivityLine } from "../components/molecules";
import * as api from "../services/api";

const filters = [
  { key: "all", label: "All" },
  { key: "navigate", label: "Navigate" },
  { key: "snapshot", label: "Snapshot" },
  { key: "action", label: "Actions" },
];

export default function AgentsPage() {
  const {
    agents,
    selectedAgentId,
    events,
    eventFilter,
    setAgents,
    setSelectedAgentId,
    setEventFilter,
  } = useAppStore();

  const loadAgents = async () => {
    try {
      const data = await api.fetchAgents();
      setAgents(data);
    } catch (e) {
      console.error("Failed to load agents", e);
    }
  };

  // Agents are loaded via SSE init event — only load if empty
  // Load once on mount if empty — intentionally omitting deps to avoid refetch loops
  useEffect(() => {
    if (agents.length === 0) {
      loadAgents();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const filteredEvents = events.filter((e) => {
    if (selectedAgentId && e.agentId !== selectedAgentId) return false;
    if (eventFilter === "all") return true;
    if (eventFilter === "action") return e.type === "action";
    return e.type === eventFilter;
  });

  return (
    <div className="flex flex-1 flex-col">
      {/* Agents carousel */}
      <div className="shrink-0 border-b border-border-subtle bg-bg-surface">
        <div className="flex items-center gap-2 overflow-x-auto px-4 py-3">
          <button
            className={`shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-all ${
              !selectedAgentId
                ? "bg-primary text-white"
                : "bg-bg-elevated text-text-secondary hover:bg-bg-elevated/80"
            }`}
            onClick={() => setSelectedAgentId(null)}
          >
            All
          </button>
          {agents.length === 0 ? (
            <span className="shrink-0 text-sm text-text-muted">
              🦀 No agents yet — make an API call with{" "}
              <code className="text-primary">X-Agent-Id</code>
            </span>
          ) : (
            agents.map((agent) => (
              <button
                key={agent.id}
                className={`shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-all ${
                  selectedAgentId === agent.id
                    ? "bg-primary text-white"
                    : "bg-bg-elevated text-text-secondary hover:bg-bg-elevated/80"
                }`}
                onClick={() => setSelectedAgentId(agent.id)}
                title={`${agent.requestCount} requests`}
              >
                {agent.name || agent.id}
                <span className="ml-1.5 text-xs opacity-70">
                  ({agent.requestCount})
                </span>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Activity feed */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <div className="flex items-center justify-between border-b border-border-subtle bg-bg-surface px-4 py-2">
          <h2 className="text-sm font-semibold text-text-secondary">
            Activity Feed
          </h2>
          <div className="flex gap-1">
            {filters.map((f) => (
              <button
                key={f.key}
                className={`rounded px-2 py-1 text-xs font-medium transition-all ${
                  eventFilter === f.key
                    ? "bg-primary/10 text-primary"
                    : "text-text-muted hover:bg-bg-elevated hover:text-text-secondary"
                }`}
                onClick={() => setEventFilter(f.key)}
              >
                {f.label}
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1 overflow-auto">
          {filteredEvents.length === 0 ? (
            <EmptyState title="Waiting for events..." icon="📡" />
          ) : (
            filteredEvents.map((event) => (
              <ActivityLine key={event.id} event={event} />
            ))
          )}
        </div>
      </div>
    </div>
  );
}
