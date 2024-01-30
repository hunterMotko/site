package tools

type Tool struct {
  Name string
  Description string
  Link string
}

type Tools struct {
  Terminal []Tool
  Neovim []Tool
}

func newTerm() []Tool {
  var tools = []Tool{
    {
      Name: "bat",
      Description: "A cat(1) clone with wings.",
      Link: "https://github.com/sharkdp/bat",
    },
    {
      Name: "fzf",
      Description: "A command-line fuzzy finder.",
      Link: "https://github.com/junegunn/fzf",
    },
    {
      Name: "gtop",
      Description: "System monitoring dashboard for terminal.",
      Link: "https://github.com/aksakalli/gtop",
    },
    {
      Name: "httpie",
      Description: "Modern command line HTTP client.",
      Link: "https://httpie.io/",
    },
    {
      Name: "ripgrep",
      Description: "Recursively searches directories for a regex pattern.",
      Link: "https://github.com/BurntSushi/ripgrep",
    },
    {
      Name: "stow",
      Description: "Organize your dotfiles neatly.",
      Link: "https://www.gnu.org/software/stow/",
    },
    {
      Name: "tldr",
      Description: "to long; didn't read for the command line.",
      Link: "https://github.com/tldr-pages/tldr",
    },
    {
      Name: "tmux",
      Description: "A terminal multiplexer.",
      Link: "https://github.com/tmux/tmux/wiki",
    },
  }
  return tools
}

func newNvim() []Tool {
  return []Tool{
    {
      Name: "fugitive",
      Description: "A Git wrapper so awesome, it should be illegal.",
      Link: "https://github.com/tpope/vim-fugitive",
    },
    {
      Name: "harpoon",
      Description: "A navigation and window manager for Neovim.",
      Link: "https://github.com/ThePrimeagen/harpoon",
    },
    {
      Name: "lsp-zero",
      Description: "A library to make the configuration of Neovim's built-in language client easier.",
      Link: "https://github.com/VonHeikemen/lsp-zero.nvim",
    },
    {
      Name: "telescope",
      Description: "A highly extendable fuzzy finder over lists.",
      Link: "https://github.com/nvim-telescope/telescope.nvim",
    },
    {
      Name: "tree-sitter",
      Description: "A parser generator tool and an incremental parsing library.",
      Link: "https://github.com/nvim-treesitter/nvim-treesitter",
    },
    {
      Name: "undotree",
      Description: "A simple undo history visualizer for Neovim.",
      Link: "https://github.com/mbbill/undotree",
    },
  }
}

func GetTools() Tools {
  return Tools{
    Terminal: newTerm(),
    Neovim: newNvim(),
  }
}
