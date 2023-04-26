use project::{Configuration, Identifier, Project};
use walkdir::WalkDir;

use std::{env::args, path::Path};

mod project;

fn main() {
    let args: Vec<String> = args().collect();

    let default_path = String::from(".");
    let path = args.get(1).unwrap_or(&default_path);
    let path = Path::new(path);

    let projects = vec![
        Project::new(
            "Javascript".to_string(),
            vec![Configuration::new(
                Identifier::File("package.json".to_string()),
                vec!["node_modules".to_string()],
            )],
        ),
        Project::new(
            "PHP".to_string(),
            vec![Configuration::new(
                Identifier::File("composer.json".to_string()),
                vec!["vendor".to_string()],
            )],
        ),
        Project::new(
            "Python".to_string(),
            vec![
                Configuration::new(
                    Identifier::FileExtensionInDirectory(
                        "__pycache__".to_string(),
                        "pyc".to_string(),
                    ),
                    vec!["__pycache__".to_string()],
                ),
                Configuration::new(
                    Identifier::File("pyvenv.cfg".to_string()),
                    vec!["".to_string()],
                ),
            ],
        ),
        Project::new(
            "Rust".to_string(),
            vec![Configuration::new(
                Identifier::File("Cargo.toml".to_string()),
                vec!["target".to_string()],
            )],
        ),
        Project::new(
            "Swift".to_string(),
            vec![Configuration::new(
                Identifier::File("Package.swift".to_string()),
                vec![".build".to_string()],
            )],
        ),
    ];

    let mut last_matched_target: Option<&Path> = None;
    // TODO: When a node project is found, subdirectories of the node_modules directory should be ignored
    for entry in WalkDir::new(path).into_iter().filter_map(|e| e.ok()) {
        let possible_identifier_path = entry.path();

        
        for project in &projects {
            for config in &project.configurations {
                if config.matches(possible_identifier_path) {
                    let identifier_root_path = possible_identifier_path.parent().unwrap();
                    for directory in &config.relative_targets {
                        let target_path = identifier_root_path.join(directory);
                        let target_path_clone = target_path.clone();
                        let exists = target_path_clone.exists();
                        let last_matched_target_clone = last_matched_target.clone();
                        let target_is_subdirectory_of_last_matched_target = target_path_clone.starts_with(last_matched_target_clone.unwrap_or(&Path::new("")));
                        let path_to_set = target_path_clone.as_path();
                        if exists  && !target_is_subdirectory_of_last_matched_target {
                            println!("Found {} Project at {}.", project.name, identifier_root_path.display());
                            last_matched_target = Some(path_to_set);
                        }
                    }
                }
            }
        }
    }
}
