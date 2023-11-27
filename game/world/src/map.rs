use std::collections::HashMap;

use log::debug;
use quadtree_rs::{area::AreaBuilder, point::Point, Quadtree};

use crate::entity::{Entity, EntityId, Position};

pub struct Map {
    id: String,
    quadtree: Quadtree<u16, EntityId>,
    entity_to_space_id: HashMap<EntityId, u64>,
}

impl Default for Map {
    fn default() -> Self {
        Self {
            id: String::from("classic"),
            quadtree: Quadtree::<u16, EntityId>::new(4),
            entity_to_space_id: HashMap::new(),
        }
    }
}

impl Map {
    pub fn id(&self) -> &str {
        &self.id
    }

    pub fn set_id(&mut self, id: String) -> &mut Self {
        self.id = id;

        self
    }

    pub fn load(&mut self) {
        debug!("Loading map: {}", self.id);
    }

    pub fn move_entity(&mut self, entity: Box<&dyn Entity>, position: Position) -> bool {
        let entity_space = AreaBuilder::default()
            .anchor(Point {
                x: position.0,
                y: position.1,
            })
            .dimensions(entity.base().size())
            .build()
            .unwrap();
        let space_id = self
            .quadtree
            .insert(entity_space, entity.base().id().clone());

        match space_id {
            Some(id) => {
                if let Some(old_id) = self.entity_to_space_id.get_mut(entity.base().id()) {
                    self.quadtree.delete_by_handle(*old_id);
                    *old_id = id;
                } else {
                    self.entity_to_space_id
                        .insert(entity.base().id().clone(), id);
                }

                true
            }
            _ => false,
        }
    }
}
