package core

import (
	"sort"
)

type HandlerFunc func(ctx *Ptz) error

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Category    string
	OwnerOnly   bool
	GroupOnly   bool
	AdminOnly   bool
	BotAdmin    bool
	Handler     HandlerFunc
}

type Registry struct {
	commands map[string]*Command
}

func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]*Command)}
}

func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

func (r *Registry) Get(name string) (*Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *Registry) All() map[string]*Command {
	return r.commands
}

func (r *Registry) ByCategory() map[string][]*Command {
	seen := map[string]struct{}{}
	result := map[string][]*Command{}
	for _, cmd := range r.commands {
		if _, ok := seen[cmd.Name]; ok {
			continue
		}
		seen[cmd.Name] = struct{}{}
		result[cmd.Category] = append(result[cmd.Category], cmd)
	}
	for cat := range result {
		sort.Slice(result[cat], func(i, j int) bool {
			return result[cat][i].Name < result[cat][j].Name
		})
	}
	return result
}

func (r *Registry) Categories() []string {
	cats := map[string]struct{}{}
	for _, cmd := range r.commands {
		cats[cmd.Category] = struct{}{}
	}
	result := make([]string, 0, len(cats))
	for c := range cats {
		result = append(result, c)
	}
	sort.Strings(result)
	return result
}

var globalRegistry = NewRegistry()

func GlobalRegistry() *Registry {
	return globalRegistry
}

func Use(cmd *Command) {
	globalRegistry.Register(cmd)
}

func (c *Command) Execute(ptz *Ptz) error {
	if c.GroupOnly && !ptz.IsGroup {
		return nil
	}

	if c.OwnerOnly && !ptz.IsOwner() {
		return nil
	}

	if c.AdminOnly || c.BotAdmin {
		if err := ptz.LoadGroupInfo(); err != nil {
			ptz.Bot.Log.Warnf("LoadGroupInfo error for %s: %v", ptz.Chat, err)
			return nil
		}
	}

	if c.AdminOnly && !ptz.IsAdmin() && !ptz.IsOwner() {
		return nil
	}

	if c.BotAdmin && !ptz.IsBotAdmin() {
		ptz.ReplyText("❌ Bot harus jadi admin dulu.")
		return nil
	}

	return c.Handler(ptz)
}
