use std::collections::HashMap;

use crate::{action::ActionType, entity::player::Player, map::Map};

pub struct World {
    is_started: bool,
    map: Map,
    players: HashMap<String, Player>,
}

impl Default for World {
    fn default() -> Self {
        Self {
            is_started: false,
            map: Map::default(),
            players: HashMap::new(),
        }
    }
}

impl World {
    pub fn map(&mut self, id: String) -> &mut Self {
        self.map.set_id(id);

        self
    }

    pub fn players(&mut self, ids: Vec<String>) -> &mut Self {
        for id in ids {
            self.players
                .insert(id.clone(), Player::new(id, (64, 64), (0, 0)));
        }

        self
    }

    pub fn build(&mut self) {
        self.map.load()
    }

    pub fn start(&mut self) {
        if self.is_started {
            return;
        }

        self.is_started = true
    }

    pub fn control_player(&mut self, id: &str, action: ActionType) {
        let Some(player) = self.players.get_mut(id) else {
            return;
        };

        match action {
            ActionType::Attack => player.attack(),
            ActionType::Move(position) => {
                if self.map.move_entity(Box::new(player), position) {
                    player.moves(position);
                }
            }
        }
    }
}
